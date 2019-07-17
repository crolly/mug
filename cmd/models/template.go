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
	Resources map[string]SAMFunction `yaml:"Resources"`
}

// SAMFunction ...
type SAMFunction struct {
	Type       string    `yaml:"Type"`
	Properties SAMFnProp `yaml:"Properties"`
}

// SAMFnProp ...
type SAMFnProp struct {
	Runtime string              `yaml:"Runtime"`
	Handler string              `yaml:"Handler"`
	CodeURI string              `yaml:"CodeUri"`
	Events  map[string]SAMEvent `yaml:"Events"`
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

// NewTemplate returns a new TemplateConfig
func NewTemplate() TemplateConfig {
	return TemplateConfig{
		Transform: "AWS::Serverless-2016-10-31",
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
							Path:   "/" + f.Events[0].HTTP.Path,
							Method: f.Events[0].HTTP.Method,
						},
					},
				},
			},
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
