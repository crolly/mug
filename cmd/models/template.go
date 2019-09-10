package models

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	Runtime     string              `yaml:"Runtime,omitempty"`
	Handler     string              `yaml:"Handler,omitempty"`
	CodeURI     string              `yaml:"CodeUri,omitempty"`
	Events      map[string]SAMEvent `yaml:"Events,omitempty"`
	Environment FnEnvironment       `yaml:"Environment,omitempty"`
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

	// add environments
	for key, val := range s.Provider.Environments {
		// get value in case it's stored in secret
		if strings.HasPrefix(val, "${file") {
			re := regexp.MustCompile(`\$\{file\((.*?)\):(.*?)\}`)
			reFound := re.FindAllStringSubmatch(val, 3)[0]
			fileName := reFound[1]
			envKey := reFound[2]
			f, err := ioutil.ReadFile(filepath.Join(s.ProjectPath, "functions", r, fileName))
			if err != nil {
				panic(err.Error())
			}

			envs := map[string]string{}
			err = yaml.Unmarshal(f, envs)
			if err != nil {
				panic(err.Error())
			}

			t.Globals.Function.Environment.Variables[envKey] = envs[envKey]
		} else {
			t.Globals.Function.Environment.Variables[key] = val
		}
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
