package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

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

// ResourceConfig represents mu's configuration for resources
type ResourceConfig struct {
	ProjectName string                 `json:"projectName"`
	ProjectPath string                 `json:"projectPath"`
	ImportPath  string                 `json:"importPath"`
	Region      string                 `json:"region"`
	Resources   map[string]*Resource   `json:"resources"`
	Functions   map[string][]*Function `json:"functions"`

	Authentication bool `json:"authentication"`

	Environments map[string]string `json:"-"`
}

// Resource represents a single Resource of the project's config
type Resource struct {
	Ident         flect.Ident                    `json:"ident"`
	Attributes    map[string]AttributeDefinition `json:"attributes"`
	KeySchema     map[string]string              `json:"key_schema"`
	CompositeKey  bool                           `json:"composite_key"`
	BillingMode   string                         `json:"billing_mode"`
	CapacityUnits map[string]byte                `json:"capacity_units"`
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

// ReadConfig return the ResourceConfig from the working directory
func ReadConfig() ResourceConfig {
	wd := GetWorkingDir()

	configFile, err := os.Open(filepath.Join(wd, "mug.config.json"))
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config ResourceConfig

	json.Unmarshal(data, &config)

	// make sure map exists
	if len(config.Resources) == 0 {
		config.Resources = make(map[string]*Resource)
	}

	return config
}

// Write method to write the config back to disk
func (c *ResourceConfig) Write() {
	fileName := filepath.Join(c.ProjectPath, "mug.config.json")

	configJSON, _ := json.MarshalIndent(c, "", "  ")
	_ = ioutil.WriteFile(fileName, configJSON, 0644)
}

// AddFunction adds a given function to the given resource name of the configuration
func (c *ResourceConfig) AddFunction(resourceName string, functionName string, path string, method string) (string, *Function) {
	if resourceName == "" {
		resourceName = "_"
	}

	ident := flect.New(resourceName)
	fName := getFuncName(ident, functionName)

	f := &Function{
		Name:    fName,
		Handler: functionName,
		Path:    path,
		Method:  method,
	}

	rCamel := ident.Camelize().String()
	c.Functions[rCamel] = append(c.Functions[rCamel], f)

	return rCamel, f
}

// RemoveFunction removes a given function from the given resource name of the configuration
func (c *ResourceConfig) RemoveFunction(resourceName string, functionName string) {
	if resourceName == "" {
		resourceName = "_"
	}

	ident := flect.New(resourceName)
	rCamel := ident.Camelize().String()
	name := getFuncName(ident, functionName)

	for i, f := range c.Functions[rCamel] {
		if name == f.Name {
			c.Functions[rCamel] = append(c.Functions[rCamel][:i], c.Functions[rCamel][i+1:]...)

			return
		}
	}
}

// RemoveResource removes a given resource from the configuration
func (c *ResourceConfig) RemoveResource(resourceName string) {
	delete(c.Resources, resourceName)
	delete(c.Functions, resourceName)
}

// getFuncName returns the generated function name for a given resource ident and a functionName
func getFuncName(ident flect.Ident, functionName string) string {
	if ident.String() == "_" {
		return functionName
	}

	return functionName + "_" + ident.Singularize().String()

}

// GetServerlessConfig returns the ServerlessConfig (serverless.yml contents) for a given resource/function group name
func (m MUGConfig) GetServerlessConfig(n string) ServerlessConfig {
	s, err := os.Open(filepath.Join(m.ProjectPath, "functions", n, "serverless.yml"))
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	data, err := ioutil.ReadAll(s)
	if err != nil {
		log.Fatal(err)
	}

	var c ServerlessConfig

	yaml.Unmarshal(data, &c)

	return c
}

// ReadMUGConfig ...
func ReadMUGConfig() MUGConfig {
	data := readDataFromFile(filepath.Join(GetWorkingDir(), "mug.config.json"))

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
func (m MUGConfig) NewServerlessConfig() ServerlessConfig {
	s := NewDefaultServerlessConfig()
	s.Service = m.ProjectName
	s.Provider.Region = m.Region

	return s
}
