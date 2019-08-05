package models

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

// ServerlessConfig ...
type ServerlessConfig struct {
	Service   Service
	Provider  Provider
	Package   Package
	Functions map[string]*ServerlessFunction `yaml:",omitempty"`
	Layers    map[string]Layer               `yaml:",omitempty"`
	Resources Resources                      `yaml:",omitempty"`
}

// Service ...
type Service struct {
	Name         string
	AWSKMSKeyARN string `yaml:"awsKmsKeyArn,omitempty"`
}

// Provider ...
type Provider struct {
	Name                string
	Runtime             string
	Region              string
	Stage               string
	StackName           string            `yaml:"stackName,omitempty"`
	APIName             string            `yaml:"apiName,omitempty"`
	Profile             string            `yaml:",omitempty"`
	MemSize             string            `yaml:"memorySize,omitempty"`
	ReservedConcurrency int               `yaml:"reservedConcurrency,omitempty"`
	Timeout             int               `yaml:"timeout,omitempty"`
	logRetention        int               `yaml:"logRetentionInDays,omitempty"`
	DeyploymentBucket   DeploymentBucket  `yaml:"deploymendBucket,omitempty"`
	DeploymentPrefix    string            `yaml:"deploymentPrefix,omitempty"`
	Environments        map[string]string `yaml:"environment,omitempty"`
	RoleStatements      []RoleStatement   `yaml:"iamRoleStatements"`
	EndpointType        string            `yaml:"endpointType,omitempty"`
	NotificationARNs    []string          `yaml:"notificationArns,omitempty"`
	Tags                map[string]string `yaml:",omitempty"`
	Tracing             TracingConfig     `yaml:",omitempty"`
	Logs                LogConfig         `yaml:",omitempty"`
}

// DeploymentBucket ...
type DeploymentBucket struct {
	Name string
	SSE  string            `yaml:"serverSideEncryption,omitempty"`
	Tags map[string]string `yaml:",omitempty"`
}

// RoleStatement ...
type RoleStatement struct {
	Effect   string   `yaml:"Effect"`
	Actions  []string `yaml:"Action"`
	Resource string   `yaml:"Resource"`
}

// TracingConfig ...
type TracingConfig struct {
	APIGateway bool `yaml:"apiGateway,omitempty"`
	Lambda     bool `yaml:",omitempty"`
}

// LogConfig ...
type LogConfig struct {
	RESTAPI   bool `yaml:"restApi,omitempty"`
	WebSocket bool `yaml:"websocket,omitempty"`
}

// Package ...
type Package struct {
	Excludes     []string `yaml:"exclude,omitempty"`
	Includes     []string `yaml:"include,omitempty"`
	Individually bool     `yaml:",omitempty"`
}

// ServerlessFunction ...
type ServerlessFunction struct {
	Handler             string
	Package             Package           `yaml:",omitempty"`
	Name                string            `yaml:",omitempty"`
	Description         string            `yaml:",omitempty"`
	MemorySize          int               `yaml:",omitempty"`
	ReservedConcurrency int               `yaml:"reservedConcurrency,omitempty"`
	RunTime             string            `yaml:"runtime,omitempty"`
	Timeout             int               `yaml:",omitempty"`
	AWSKMSKeyARN        string            `yaml:"awsKmsKeyArn,omitempty"`
	Environments        map[string]string `yaml:"environment,omitempty"`
	Tags                map[string]string `yaml:",omitempty"`
	Layers              []string          `yaml:",omitempty"`
	Events              []Events
}

// Events ...
type Events struct {
	HTTP            *HTTPEvent      `yaml:"http,omitempty"`
	WebSocket       *WebSocketEvent `yaml:"websocket,omitempty"`
	S3              *S3Event        `yaml:"s3,omitempty"`
	Schedule        *ScheduleEvent  `yaml:"schedule,omitempty"`
	SNS             *SNSEvent       `yaml:"sns,omitempty"`
	SQS             *SQSEvent       `yaml:"sqs,omitempty"`
	Stream          *StreamEvent    `yaml:"stream,omitempty"`
	AlexaSkill      *AlexaEvent     `yaml:"alexaSkill,omitempty"`
	AlexaSmartHome  *AlexaEvent     `yaml:"alexaSmartHome,omitempty"`
	IOT             *IOTEvent       `yaml:"iot,omitempty"`
	CognitoUserPool *CognitoEvent   `yaml:"cognitoUserPool,omitempty"`
	ALB             *ALBEvent       `yaml:"alb,omitempty"`
}

// HTTPEvent ...
type HTTPEvent struct {
	Path       string
	Method     string
	CORS       bool        `yaml:",omitempty"`
	Private    bool        `yaml:",omitempty"`
	Authorizer *Authorizer `yaml:",omitempty"`
	Scopes     []string    `yaml:",omitempty"`
}

// WebSocketEvent ...
type WebSocketEvent struct {
	Route      string
	Authorizer *Authorizer `yaml:",omitempty"`
}

// S3Event ...
type S3Event struct {
	Bucket string
	Event  string
	Rules  map[string]string `yaml:",omitempty"`
}

// ScheduleEvent ...
type ScheduleEvent struct {
	Name        string
	Description string `yaml:",omitempty"`
	Rate        string
	Enabled     bool   `yaml:",omitempty"`
	InputPath   string `yaml:"inputPath,omitempty"`
}

// SNSEvent ...
type SNSEvent struct {
	TopicName   string `yaml:"topicName"`
	DisplayName string `yaml:"displayName,omitempty"`
}

// SQSEvent ...
type SQSEvent struct {
	ARN       string `yaml:"arn"`
	BatchSize int    `yaml:"batchSize,omitempty"`
}

// StreamEvent ...
type StreamEvent struct {
	ARN              string `yaml:"arn"`
	BatchSize        int    `yaml:"batchSize,omitempty"`
	StartingPosition string `yaml:"startingPosition,omitempty"`
	Enabled          bool   `yaml:",omitempty"`
}

// AlexaEvent ...
type AlexaEvent struct {
	AppID   string `yaml:"appId"`
	Enabled bool   `yaml:",omitempty"`
}

// IOTEvent ...
type IOTEvent struct {
	Name        string `yaml:",omitempty"`
	Description string `yaml:",omitempty"`
	Enabled     bool   `yaml:",omitempty"`
	SQL         string `yaml:"sql,omitempty"`
	SQLVersion  string `yaml:"sqlVersion,omitempty"`
}

// CognitoEvent ...
type CognitoEvent struct {
	Pool     string
	Trigger  string
	Existing bool
}

// ALBEvent ...
type ALBEvent struct {
	ListenerARN string            `yaml:"listenerArn"`
	Priority    int               `yaml:",omitempty"`
	Conditions  map[string]string `yaml:",omitempty"`
}

// Authorizer ...
type Authorizer struct {
	ARN                          string
	Name                         string `yaml:",omitempty"`
	ResultTTL                    int    `yaml:"resultTtlInSeconds,omitempty"`
	IdentitySource               string `yaml:"identitySource,omitempty"`
	IdentityValidationExpression string `yaml:"identityValidationExpression,omitempty"`
	Type                         string `yaml:",omitempty"`
}

// Layer ...
type Layer struct {
	Path               string
	Name               string   `yaml:",omitempty"`
	Description        string   `yaml:",omitempty"`
	CompatibleRuntimes []string `yaml:"compatibleRuntimes,omitempty"`
	License            string   `yaml:"licenseInfo,omitempty"`
	AllowedAccounts    []string `yaml:"allowedAccounts,omitempty"`
	Retain             bool     `yaml:",omitempty"`
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
	AttributeDefinitions   []AttributeDef         `yaml:"AttributeDefinitions,omitempty"`
	KeySchema              []KeySchema            `yaml:"KeySchema,omitempty"`
	ProvisionedThroughput  *ProvisionedThroughput `yaml:"ProvisionedThroughput,omitempty"`
	BillingMode            string                 `yaml:"BillingMode,omitempty"`
	TableName              string                 `yaml:"TableName,omitempty"`
	LocalSecondaryIndexes  []LocalIndex           `yaml:"LocalSecondaryIndexes,omitempty"`
	GlobalSecondaryIndexes []GlobalIndex          `yaml:"GlobalSecondaryIndexes,omitempty"`
	TTLSpecification       TTLSpecification       `yaml:"TimeToLiveSpecification,omitempty"`

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

// LocalIndex ...
type LocalIndex struct {
	IndexName  string      `yaml:"IndexName"`
	KeySchema  []KeySchema `yaml:"KeySchema"`
	Projection Projection  `yaml:"Projection"`
}

// GlobalIndex ...
type GlobalIndex struct {
	IndexName             string                 `yaml:"IndexName"`
	KeySchema             []KeySchema            `yaml:"KeySchema"`
	Projection            Projection             `yaml:"Projection"`
	ProvisionedThroughput *ProvisionedThroughput `yaml:"ProvisionedTroughput,omitempty"`
}

// Projection ...
type Projection struct {
	NonKeyAttributes []string `yaml:"NonKeyAttributes,omitempty"`
	ProjectionType   string   `yaml:"ProjectionType,omitempty"`
}

// TTLSpecification ...
type TTLSpecification struct {
	AttributeName string `yaml:"AttributeName"`
	Enabled       bool   `yaml:"Enabled"`
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
			Individually: true,
		},
	}

	return s
}

func (s *ServerlessConfig) updateSecrets(path string) {
	// read secrets
	secrets := map[string]string{}
	data, err := readDataFromFile(filepath.Join(path, "secrets.yml"))
	if err != nil && !os.IsNotExist(err) {
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
		s.Provider.Environments[k] = "${file(secrets.yml):" + k + "}"
	}
}

func (s *ServerlessConfig) updateEnv(projectPath, resourcePath string) {
	// read env file
	if env, _ := godotenv.Read(filepath.Join(projectPath, ".env"), filepath.Join(resourcePath, ".env")); len(env) > 0 {
		// merge environment into ServerlessConfig
		if err := mergo.Merge(&s.Provider.Environments, env); err != nil {
			log.Fatal(err)
		}
	}
}

// Write writes the ServerlessConfig to serverless.yml for the given project path and model name
func (s *ServerlessConfig) Write(projectPath, mn string) {
	// resource path
	rp := filepath.Join(projectPath, "functions", mn)

	// update secrets
	s.updateSecrets(rp)
	s.updateEnv(projectPath, rp)

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
func (s *ServerlessConfig) SetResourceWithModel(r *NewResource, m Model, projectName string) {
	tableName := projectName + "-" + r.Ident.Pluralize().String() + "-${opt:stage, self:provider.stage}"
	rd := &ResourceDefinition{
		Type:           "AWS::DynamoDB::Table",
		DeletionPolicy: "Retain",
		Properties: Properties{
			TableName: tableName,
		},
	}

	// set key attributes and key schema
	for _, a := range r.Attributes {
		rd.Properties.AttributeDefinitions = append(rd.Properties.AttributeDefinitions, AttributeDef{
			AttributeName: a.Ident.String(),
			AttributeType: a.AwsType,
		})
	}

	rd.Properties.KeySchema = []KeySchema{
		{
			AttributeName: m.KeySchema["HASH"],
			KeyType:       "HASH",
		},
	}

	if m.CompositeKey {
		rd.Properties.KeySchema = append(rd.Properties.KeySchema, KeySchema{
			AttributeName: m.KeySchema["RANGE"],
			KeyType:       "RANGE",
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
		r.Ident.ToUpper().String() + "_TABLE_NAME": tableName,
	}
}

// SetFunctions sets a slice of Functions to the ServerlessConfig
func (s *ServerlessConfig) SetFunctions(fns []*Function) {
	s.Functions = map[string]*ServerlessFunction{}

	for _, fn := range fns {
		s.AddFunction(fn)
	}
}

// AddFunction adds a function to the ServerlessConfig
func (s *ServerlessConfig) AddFunction(fn *Function) {
	// make sure map exists
	if len(s.Functions) == 0 {
		s.Functions = map[string]*ServerlessFunction{}
	}
	s.Functions[fn.Name] = &ServerlessFunction{
		Handler: "bin/" + fn.Handler,
		Package: Package{
			Includes: []string{
				"bin/" + fn.Handler,
			},
			Excludes: []string{
				"./**",
			},
		},
		Events: []Events{
			Events{
				HTTP: &HTTPEvent{
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

	for _, fn := range s.Functions {
		if !Contains(excludeSlice, strings.TrimPrefix(fn.Handler, "bin/")) {
			fn.addAuth()
		}
	}
}

// RemoveAuth removes Authorization from the ServerlessConfig
func (s *ServerlessConfig) RemoveAuth() {
	for _, fn := range s.Functions {
		fn.removeAuth()
	}
}

// addAuth adds the authorizer reference to the ServerlessFunction
func (f *ServerlessFunction) addAuth() {
	f.Events[0].HTTP.Authorizer = &Authorizer{
		ARN: "${file(secrets.yml):COGNITO_USER_POOL}",
	}
}

// removeAuth removes the authorizer reference to the ServerlessFunction
func (f *ServerlessFunction) removeAuth() {
	f.Events[0].HTTP.Authorizer = nil
}
