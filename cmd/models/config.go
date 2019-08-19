package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gobuffalo/flect"
)

// MUGConfig ...
type MUGConfig struct {
	ProjectName string
	ProjectPath string
	ImportPath  string
	Region      string
	Resources   map[string]*NewResource
}

// NewResource ...
type NewResource struct {
	Ident      flect.Ident
	Attributes map[string]AttributeDefinition
}

// AttributeDefinition represents the definition of a resource's attribute
type AttributeDefinition struct {
	Ident   flect.Ident `json:"ident"`
	AwsType string      `json:"aws_type"`
}

// Function represents a Function
type Function struct {
	Name    string `json:"name"`
	Handler string `json:"handler"`
	Path    string `json:"path"`
	Method  string `json:"method"`

	Authentication bool `json:"authentication"`
}

// ReadMUGConfig ...
func ReadMUGConfig() MUGConfig {
	data, err := readDataFromFile(filepath.Join(GetWorkingDir(), "mug.config.json"))
	if err != nil {
		log.Fatal(err)
	}

	var config MUGConfig
	json.Unmarshal(data, &config)

	// make sure map exists
	if len(config.Resources) == 0 {
		config.Resources = make(map[string]*NewResource)
	}

	return config
}

// Write write the MUGConfig to mug.config.json in the project path
func (m MUGConfig) Write() {
	f := filepath.Join(m.ProjectPath, "mug.config.json")

	json, _ := json.MarshalIndent(m, "", "  ")
	err := ioutil.WriteFile(f, json, 0644)

	if err != nil {
		log.Fatal(err)
	}
}

// NewServerlessConfig return a new ServerlessConfig with the attributes from the MUGConfig
// NewFromResourceConfig returns a ServerlessConfig from a provided ResourceConfig
func (m MUGConfig) NewServerlessConfig(resource string) ServerlessConfig {
	s := NewDefaultServerlessConfig()
	s.Service = Service{Name: m.ProjectName + "-" + resource}
	s.Provider.Region = m.Region

	return s
}

// ReadServerlessConfig reads the ServerlessConfig from serverless.yml in the resource or function group directory.
// If a serverless.yml file does not exist, a new default ServerlessConfig is returned
func (m MUGConfig) ReadServerlessConfig(rn string) ServerlessConfig {
	var sc ServerlessConfig
	data, err := readDataFromFile(filepath.Join(m.ProjectPath, "functions", rn, "serverless.yml"))
	if err == nil {
		if err := yaml.Unmarshal(data, &sc); err != nil {
			log.Fatal(err)
		}
	} else if os.IsNotExist(err) {
		// file doesn't exist return default ServerlessConfig
		sc = m.NewServerlessConfig(rn)

	}

	return sc
}

// RemoveResource removes a given resource from the MUGConfig and ServerlessConfig
func (m *MUGConfig) RemoveResource(rN string) {
	// remove from MUGConfig
	delete(m.Resources, rN)
}

func (m MUGConfig) clearFolder(path string) {
	if err := os.RemoveAll(filepath.Join(m.ProjectPath, "functions", path)); err != nil {
		log.Fatal(err)
	}
}

func (m MUGConfig) renderMakefile(t *template.Template, r string) {
	// open file and execute template
	f, err := os.Create(filepath.Join(m.ProjectPath, "Makefile"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sc := m.ReadServerlessConfig(r)
	// execute template and save to file
	data := map[string]interface{}{
		"Functions": sc.Functions,
		"Resource":  r,
	}

	err = t.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Makefile generated.")
}

func (m MUGConfig) make(list []string, path, target string, test bool) {
	// load Makefile template
	t := LoadTemplateFromBox(MakeBox, "Makefile.tmpl")

	for _, r := range list {
		// clear the debug binaries
		m.clearFolder(filepath.Join(r, path))
		// render for each resource/ function group
		m.renderMakefile(t, r)
		// run test if flag indicates so
		if test {
			log.Println("Run tests")
			RunCmd("make", "test")
		}
		// and run the build
		RunCmd("make", target)
	}
}

// MakeDebug renders the Makefile and builds the debug binaries
func (m MUGConfig) MakeDebug(list []string) {
	m.make(list, "debug", "debug", false)
}

// MakeBuild renders the Makefile and builds the binaries
func (m MUGConfig) MakeBuild(list []string, test bool) {
	m.make(list, "bin", "build", test)
}

// CreateResourceTables creates the tables in the local DynamoDB named by the given mode
func (m MUGConfig) CreateResourceTables(list []string, mode string, overwrite bool) {
	// create service to dynamodb
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String(m.Region),
	}))
	svc := dynamodb.New(sess)

	// get list of tables
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("Error during creation of resource tables: %s\n", err)
	}

	tables := make(map[string]bool)
	for _, t := range result.TableNames {
		tables[*t] = true
	}

	// iterate over resources
	for n, r := range m.Resources {
		if Contains(list, n) {
			sc := m.ReadServerlessConfig(n)
			rName := r.Ident.Pascalize().String() + "DynamoDbTable"
			res := sc.Resources.Resources[rName]
			if res == nil {
				log.Fatalf("Resourse %s not valid. Please check your serverless.yml or your command.", rName)
			}
			props := res.Properties
			tableName := m.ProjectName + "-" + r.Ident.Pluralize().Camelize().String() + "-" + mode

			if tables[tableName] {
				if overwrite {
					deleteTable(svc, tableName)
					createTableForResource(svc, tableName, props)
				} else {
					log.Printf("Table %s already exists, skipping creation...", tableName)
				}
			} else {
				createTableForResource(svc, tableName, props)
			}
		}
	}
}

func createTableForResource(svc *dynamodb.DynamoDB, tableName string, props Properties) {
	// create the table input for the resource
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
	}

	// get keySchema
	keySchema := []*dynamodb.KeySchemaElement{}
	for _, k := range props.KeySchema {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{
			AttributeName: aws.String(flect.New(k.AttributeName).Underscore().String()),
			KeyType:       aws.String(k.KeyType),
		})
	}

	// get attributes
	attributes := []*dynamodb.AttributeDefinition{}
	for _, a := range props.AttributeDefinitions {
		attributes = append(attributes, &dynamodb.AttributeDefinition{
			AttributeName: aws.String(flect.New(a.AttributeName).Underscore().String()),
			AttributeType: aws.String(a.AttributeType),
		})
	}

	// get throughput
	throughput := &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(1),
		WriteCapacityUnits: aws.Int64(1),
	}
	if props.ProvisionedThroughput != nil {
		throughput = &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(props.ProvisionedThroughput.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(props.ProvisionedThroughput.WriteCapacityUnits),
		}
	}

	// get the local secondary indexex
	// TODO: support for non index attributes
	lsi := []*dynamodb.LocalSecondaryIndex{}
	for _, i := range props.LocalSecondaryIndexes {
		keySchema := []*dynamodb.KeySchemaElement{}
		for _, k := range i.KeySchema {
			keySchema = append(keySchema, &dynamodb.KeySchemaElement{
				AttributeName: aws.String(flect.New(k.AttributeName).Underscore().String()),
				KeyType:       aws.String(k.KeyType),
			})
		}
		lsi = append(lsi, &dynamodb.LocalSecondaryIndex{
			IndexName: aws.String(i.IndexName),
			KeySchema: keySchema,
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String(i.Projection.ProjectionType),
			},
		})
	}

	gsi := []*dynamodb.GlobalSecondaryIndex{}
	for _, i := range props.GlobalSecondaryIndexes {
		keySchema := []*dynamodb.KeySchemaElement{}
		for _, k := range i.KeySchema {
			keySchema = append(keySchema, &dynamodb.KeySchemaElement{
				AttributeName: aws.String(flect.New(k.AttributeName).Underscore().String()),
				KeyType:       aws.String(k.KeyType),
			})
		}
		idx := &dynamodb.GlobalSecondaryIndex{
			IndexName: aws.String(i.IndexName),
			KeySchema: keySchema,
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String(i.Projection.ProjectionType),
			},
		}
		if i.ProvisionedThroughput != nil {
			idx.ProvisionedThroughput = &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(i.ProvisionedThroughput.ReadCapacityUnits),
				WriteCapacityUnits: aws.Int64(i.ProvisionedThroughput.WriteCapacityUnits),
			}
		} else {
			idx.ProvisionedThroughput = &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			}
		}
		gsi = append(gsi, idx)
	}

	// append properties to input
	if len(keySchema) > 0 {
		input.KeySchema = keySchema
	} else {
		log.Fatal("KeySchema has to be provided")
	}
	if len(attributes) == len(keySchema)+len(lsi)+len(gsi) {
		input.AttributeDefinitions = attributes
	} else {
		log.Fatal("Number of attributes defined invalid. Did you add your Local Secondary Index to the Attribute Definition?")
	}
	if len(lsi) > 0 {
		input.LocalSecondaryIndexes = lsi
	}
	if len(gsi) > 0 {
		input.GlobalSecondaryIndexes = gsi
	}
	if throughput != nil {
		input.ProvisionedThroughput = throughput
	}

	out, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Error creating table %s: %s", tableName, err)
	}

	log.Printf("Table %s created: %s", tableName, out)
}

func deleteTable(svc *dynamodb.DynamoDB, tableName string) {
	_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		log.Fatalf("Error deleting table %s: %s", tableName, err)
	}
}
