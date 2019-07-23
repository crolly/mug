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
func (m MUGConfig) NewServerlessConfig() ServerlessConfig {
	s := NewDefaultServerlessConfig()
	s.Service = Service{Name: m.ProjectName}
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
