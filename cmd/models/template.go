package models

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/gobuffalo/flect"
)

// TemplateConfig ...
type TemplateConfig struct {
	Transform string                 `yaml:"Transform"`
	Globals   GlobalConfig           `yaml:"Globals"`
	Resources map[string]SAMFunction `yaml:"Resources"`
}

// GlobalConfig ...
type GlobalConfig struct {
	Function SAMFnProp `yaml:"Function"`
}

// SAMFunction ...
type SAMFunction struct {
	Type       string    `yaml:"Type"`
	Properties SAMFnProp `yaml:"Properties"`
}

// SAMFnProp ...
type SAMFnProp struct {
	Runtime     string              `yaml:"Runtime"`
	Handler     string              `yaml:"Handler"`
	CodeURI     string              `yaml:"CodeUri"`
	Events      map[string]SAMEvent `yaml:"Events"`
	Environment FnEnvironment       `yaml:"Environment"`
}

// SAMEvent ...
type SAMEvent struct {
	Type       string  `yaml:"Type"`
	Properties SAMProp `yaml:"Properties"`
}

// SAMProp ...
type SAMProp struct {
	Path   string `yaml:"Path"`
	Method string `yaml:"Method"`
}

// FnEnvironment ...
type FnEnvironment struct {
	Variables map[string]string `yaml:"Variables"`
}

// NewTemplate returns a new TemplateConfig
func NewTemplate() TemplateConfig {
	return TemplateConfig{
		Transform: "AWS::Serverless-2016-10-31",
		Globals: GlobalConfig{
			Function: SAMFnProp{
				Environment: FnEnvironment{
					Variables: map[string]string{
						"MODE": "debug",
					},
				},
			},
		},
	}
}

// AddFunctionsFromServerlessConfig adds the functions from the given
// ServerlessConfig and resource/ function group
func (t *TemplateConfig) AddFunctionsFromServerlessConfig(s ServerlessConfig, r string) {
	if len(t.Resources) == 0 {
		t.Resources = map[string]SAMFunction{}
	}

	for n, f := range s.Functions {
		fName := flect.New(n).Camelize().String() + "Function"
		// ensure to add only http event functions
		ev := f.Events[0].HTTP
		if ev != nil {
			t.Resources[fName] = SAMFunction{
				Type: "AWS::Serverless::Function",
				Properties: SAMFnProp{
					Runtime: "go1.x",
					Handler: strings.TrimPrefix(f.Handler, "bin/"),
					CodeURI: filepath.Join(".", "functions", r, "debug"),
					Events: map[string]SAMEvent{
						"http": SAMEvent{
							Type: "Api",
							Properties: SAMProp{
								Path:   "/" + ev.Path,
								Method: ev.Method,
							},
						},
					},
				},
			}
		}
	}
}

// Write writes the TemplateConfig to template.yml
func (t *TemplateConfig) Write(projectPath string) {
	fp := filepath.Join(projectPath, "template.yml")
	// make sure directory exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			log.Fatal(err)
		}
	}

	yml, err := yaml.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fp, yml, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
