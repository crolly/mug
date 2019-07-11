package models

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"

	"github.com/joho/godotenv"

	"gopkg.in/yaml.v2"
)

// ServerlessConfig ...
type ServerlessConfig struct {
	Service   string
	Provider  Provider
	Package   Package
	Functions map[string]ServerlessFunction `yaml:",omitempty"`
	Resources Resources                     `yaml:",omitempty"`
}

// Provider ...
type Provider struct {
	Name           string
	Runtime        string
	Region         string
	Stage          string
	Environments   map[string]string `yaml:"environment,omitempty"`
	RoleStatements []RoleStatement   `yaml:"iamRoleStatements"`
}

// RoleStatement ...
type RoleStatement struct {
	Effect   string   `yaml:"Effect"`
	Actions  []string `yaml:"Action"`
	Resource string   `yaml:"Resource"`
}

// Package ...
type Package struct {
	Excludes []string `yaml:"exclude"`
	Includes []string `yaml:"include"`
}

// ServerlessFunction ...
type ServerlessFunction struct {
	Handler string
	Events  []Events
}

// Events ...
type Events struct {
	HTTP Event `yaml:"http"`
}

// Event ...
type Event struct {
	Path       string
	Method     string
	CORS       bool       `yaml:",omitempty"`
	Authorizer Authorizer `yaml:",omitempty"`
}

// Authorizer ...
type Authorizer struct {
	Type         string
	AuthorizerID Reference `yaml:"authorizerId"`
}

// Resources ...
type Resources struct {
	Resources map[string]ResourceDefinition `yaml:"Resources"`
}

// ResourceDefinition ...
type ResourceDefinition struct {
	DependsOn      []string   `yaml:"DependsOn,omitempty"`
	Type           string     `yaml:"Type"`
	DeletionPolicy string     `yaml:"DeletionPolicy,omitempty"`
	Properties     Properties `yaml:"Properties"`
}

// Properties ...
type Properties struct {
	// Properties for Resources
	AttributeDefinitions  []AttributeDef        `yaml:"AttributeDefinitions,omitempty"`
	KeySchema             []KeySchema           `yaml:"KeySchema,omitempty"`
	ProvisionedThroughput ProvisionedThroughput `yaml:"ProvisionedThroughput,omitempty"`
	BillingMode           string                `yaml:"BillingMode,omitempty"`
	TableName             string                `yaml:"TableName,omitempty"`

	//Properties for Authorizer
	Name           string    `yaml:"Name,omitempty"`
	IdentitySource string    `yaml:"IdentitySource,omitempty"`
	RestAPIID      Reference `yaml:"RestApiId,omitempty"`
	Type           string    `yaml:"Type,omitempty"`
	ProviderARNs   []string  `yaml:"ProviderARNs,omitempty"`
}

// Reference ...
type Reference struct {
	Ref string `yaml:"Ref"`
}

// AttributeDef ...
type AttributeDef struct {
	AttributeName string `yaml:"AttributeName"`
	AttributeType string `yaml:"AttributeType"`
}

// KeySchema ...
type KeySchema struct {
	AttributeName string `yaml:"AttributeName"`
	KeyType       string `yaml:"KeyType"`
}

// ProvisionedThroughput ...
type ProvisionedThroughput struct {
	ReadCapacityUnits  byte `yaml:"ReadCapacityUnits"`
	WriteCapacityUnits byte `yaml:"WriteCapacityUnits"`
}

// NewDefaultServerlessConfig return a default ServerlessConfig object
func NewDefaultServerlessConfig() ServerlessConfig {
	s := ServerlessConfig{
		Provider: Provider{
			Name:    "aws",
			Runtime: "go1.x",
			Stage:   "${opt:stage, 'dev'}",
			RoleStatements: []RoleStatement{
				RoleStatement{
					Effect: "Allow",
					Actions: []string{
						"dynamodb:DescribeTable",
						"dynamodb:Query",
						"dynamodb:Scan",
						"dynamodb:GetItem",
						"dynamodb:PutItem",
						"dynamodb:UpdateItem",
						"dynamodb:DeleteIte",
					},
					Resource: "arn:aws:dynamodb:*:*:*",
				},
			},
		},
		Package: Package{
			Excludes: []string{
				"./**",
			},
			Includes: []string{
				"bin/**",
			},
		},
	}

	return s
}

// Write writes the ServerlessConfig to serverless.yml for the given project path and model name
func (s *ServerlessConfig) Write(pp, mn string) {
	// resource path
	rp := filepath.Join(pp, "functions", mn)

	// read env file
	env, _ := godotenv.Read(filepath.Join(pp, ".env"), filepath.Join(rp, ".env"))

	// merge environment into ServerlessConfig
	if err := mergo.Merge(&s.Provider.Environments, env); err != nil {
		log.Fatal(err)
	}

	fp := filepath.Join(rp, "serverless.yml")
	// make sure directory exists
	if _, err := os.Stat(rp); os.IsNotExist(err) {
		if err := os.MkdirAll(rp, 0755); err != nil {
			log.Fatal(err)
		}
	}

	yml, err := yaml.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fp, yml, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// SetResourceWithModel sets a Resource to the ServerlessConfig
func (s *ServerlessConfig) SetResourceWithModel(r *NewResource, m Model) {
	rd := ResourceDefinition{
		Type:           "AWS::DynamoDB::Table",
		DeletionPolicy: "Retain",
		Properties: Properties{
			TableName: r.Ident.Pluralize().String() + "-${opt:stage, self:provider.stage}",
		},
	}

	// set key attributes and key schema
	for _, a := range r.Attributes {
		rd.Properties.AttributeDefinitions = append(rd.Properties.AttributeDefinitions, AttributeDef{
			AttributeName: a.Ident.String(),
			AttributeType: a.AwsType,
		})
	}

	for t, n := range m.KeySchema {
		rd.Properties.KeySchema = append(rd.Properties.KeySchema, KeySchema{
			AttributeName: n,
			KeyType:       t,
		})
	}

	// set billing mode and capacity units
	if m.BillingMode == "provisioned" {
		rd.Properties.ProvisionedThroughput = ProvisionedThroughput{
			ReadCapacityUnits:  m.CapacityUnits["read"],
			WriteCapacityUnits: m.CapacityUnits["write"],
		}
	} else if m.BillingMode == "ondemand" {
		rd.Properties.BillingMode = "PAY_PER_REQUEST"
	}

	s.Resources = Resources{
		Resources: map[string]ResourceDefinition{
			r.Ident.Pascalize().String(): rd,
		},
	}

	// set environment
	s.Provider.Environments = map[string]string{
		r.Ident.ToUpper().String() + "_TABLE_NAME": r.Ident.Pluralize().String() + "-${opt:stage, self:provider.stage}",
	}
}

// SetFunctions sets a slice of Functions to the ServerlessConfig
func (s *ServerlessConfig) SetFunctions(fns []*Function) {
	s.Functions = map[string]ServerlessFunction{}

	for _, fn := range fns {
		s.Functions[fn.Name] = ServerlessFunction{
			Handler: fn.Handler,
			Events: []Events{
				Events{
					HTTP: Event{
						Path:   fn.Path,
						Method: fn.Method,
						CORS:   true,
					},
				},
			},
		}
	}
}

// AddFunction adds a function to the ServerlessConfig
func (s *ServerlessConfig) AddFunction(fn Function) {
	// make sure map exists
	if len(s.Functions) == 0 {
		s.Functions = map[string]ServerlessFunction{}
	}
	s.Functions[fn.Name] = ServerlessFunction{
		Handler: fn.Handler,
		Events: []Events{
			Events{
				HTTP: Event{
					Path:   fn.Path,
					Method: fn.Method,
					CORS:   true,
				},
			},
		},
	}
}
