package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// ServerlessConfig ...
type ServerlessConfig struct {
	Service   string
	Provider  Provider
	Package   Package
	Functions map[string]*ServerlessFunction `yaml:",omitempty"`
	Resources Resources                      `yaml:",omitempty"`
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
	CORS       bool        `yaml:",omitempty"`
	Authorizer *Authorizer `yaml:",omitempty"`
}

// Authorizer ...
type Authorizer struct {
	Type         string
	AuthorizerID Reference `yaml:"authorizerId"`
}

// Resources ...
type Resources struct {
	Resources map[string]*ResourceDefinition `yaml:"Resources"`
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
	AttributeDefinitions  []AttributeDef         `yaml:"AttributeDefinitions,omitempty"`
	KeySchema             []KeySchema            `yaml:"KeySchema,omitempty"`
	ProvisionedThroughput *ProvisionedThroughput `yaml:"ProvisionedThroughput,omitempty"`
	BillingMode           string                 `yaml:"BillingMode,omitempty"`
	TableName             string                 `yaml:"TableName,omitempty"`

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
	ReadCapacityUnits  int64 `yaml:"ReadCapacityUnits"`
	WriteCapacityUnits int64 `yaml:"WriteCapacityUnits"`
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
func (s *ServerlessConfig) Write(projectPath, mn string) {
	// resource path
	rp := filepath.Join(projectPath, "functions", mn)

	// read secrets
	secrets := map[string]string{}
	data, err := readDataFromFile(filepath.Join(rp, "secrets.yml"))
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(data, &secrets); err != nil {
		log.Fatal(err)
	}

	// make sure map exists
	if len(s.Provider.Environments) == 0 {
		s.Provider.Environments = map[string]string{}
	}
	for k := range secrets {
		s.Provider.Environments[k] = "${file(secrets.yml):" + k
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
	rd := &ResourceDefinition{
		Type:           "AWS::DynamoDB::Table",
		DeletionPolicy: "Retain",
		Properties: Properties{
			TableName: fmt.Sprintf("${env:%s_TABLE_NAME}", m.Ident.ToUpper().String()),
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
		rd.Properties.ProvisionedThroughput = &ProvisionedThroughput{
			ReadCapacityUnits:  m.CapacityUnits["read"],
			WriteCapacityUnits: m.CapacityUnits["write"],
		}
	} else if m.BillingMode == "ondemand" {
		rd.Properties.BillingMode = "PAY_PER_REQUEST"
	}

	s.Resources = Resources{
		Resources: map[string]*ResourceDefinition{
			r.Ident.Pascalize().String() + "DynamoDbTable": rd,
		},
	}

	// set environment
	s.Provider.Environments = map[string]string{
		r.Ident.ToUpper().String() + "_TABLE_NAME": r.Ident.Pluralize().String() + "-${opt:stage, self:provider.stage}",
	}
}

// SetFunctions sets a slice of Functions to the ServerlessConfig
func (s *ServerlessConfig) SetFunctions(fns []*Function) {
	s.Functions = map[string]*ServerlessFunction{}

	for _, fn := range fns {
		s.Functions[fn.Name] = &ServerlessFunction{
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
		s.Functions = map[string]*ServerlessFunction{}
	}
	s.Functions[fn.Name] = &ServerlessFunction{
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

// RemoveFunction removes a function from the ServerlessConfig
func (s *ServerlessConfig) RemoveFunction(n string) {
	delete(s.Functions, n)
}

// AddPoolEnv adds the given cognito user pool arn as environment in .env
func (s *ServerlessConfig) AddPoolEnv(mc MUGConfig, rName, pool string) {
	path := filepath.Join(mc.ProjectPath, "functions", rName, "secrets.yml")
	var secrets map[string]string

	data, err := readDataFromFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			secrets = make(map[string]string)
		} else {
			log.Fatal(err)
		}
	}

	if err := yaml.Unmarshal(data, &secrets); err != nil {
		log.Fatal(err)
	}

	secrets["COGNITO_USER_POOL"] = pool

	yml, err := yaml.Marshal(secrets)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(path, yml, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// AddAuth adds Authorization to the ServerlessConfig
func (s *ServerlessConfig) AddAuth(excludes string) {
	excludeSlice := strings.Split(excludes, ",")
	resRequired := false

	for _, fn := range s.Functions {
		if !Contains(excludeSlice, fn.Handler) {
			resRequired = true
			fn.addAuth()
		}
	}

	if resRequired {
		s.addAuthResource()
	}
}

// RemoveAuth removes Authorization from the ServerlessConfig
func (s *ServerlessConfig) RemoveAuth() {
	s.removeResource("ApiGatewayAuthorizer")

	for _, fn := range s.Functions {
		fn.removeAuth()
	}
}

// addAuthResource adds the Authorizer Resource to the ServerlessConfig
func (s *ServerlessConfig) addAuthResource() {
	// make sure map exists
	if len(s.Resources.Resources) == 0 {
		s.Resources.Resources = map[string]*ResourceDefinition{}
	}

	s.Resources.Resources["ApiGatewayAuthorizer"] = &ResourceDefinition{
		DependsOn: []string{"ApiGatewayRestApi"},
		Type:      "AWS::ApiGateway::Authorizer",
		Properties: Properties{
			Name:           "cognito-authorizer",
			IdentitySource: "method.request.header.Authorization",
			RestAPIID: Reference{
				Ref: "ApiGatewayRestApi",
			},
			Type:         "COGNITO_USER_POOLS",
			ProviderARNs: []string{"${file(secrets.yml):COGNITO_USER_POOL}"},
		},
	}
}

func (s *ServerlessConfig) removeResource(rN string) {
	delete(s.Resources.Resources, rN)
}

// addAuth adds the authorizer reference to the ServerlessFunction
func (f *ServerlessFunction) addAuth() {
	f.Events[0].HTTP.Authorizer = &Authorizer{
		Type: "COGNITO_USER_POOLS",
		AuthorizerID: Reference{
			Ref: "ApiGatewayAuthorizer",
		},
	}
}

// removeAuth removes the authorizer reference to the ServerlessFunction
func (f *ServerlessFunction) removeAuth() {
	f.Events[0].HTTP.Authorizer = nil
}
