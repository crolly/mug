package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

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
	fName := GetFuncName(resourceName, functionName)

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
	name := GetFuncName(resourceName, functionName)

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
func (m MUGConfig) NewServerlessConfig() ServerlessConfig {
	s := NewDefaultServerlessConfig()
	s.Service = m.ProjectName
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
		sc = m.NewServerlessConfig()

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

func (m MUGConfig) make(list []string, path, target string) {
	// load Makefile template
	t := LoadTemplateFromBox(MakeBox, "Makefile.tmpl")

	for _, r := range list {
		// clear the debug binaries
		m.clearFolder(filepath.Join(r, path))
		// render for each resource/ function group
		m.renderMakefile(t, r)
		// and run the build
		RunCmd("make", target)
	}
}

// MakeDebug renders the Makefile and builds the debug binaries
func (m MUGConfig) MakeDebug(list []string) {
	m.make(list, "debug", "debug")
}

// MakeBuild renders the Makefile and builds the binaries
func (m MUGConfig) MakeBuild(list []string) {
	m.make(list, "bin", "build")
}