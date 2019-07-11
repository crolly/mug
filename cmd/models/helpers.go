package models

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/packr/v2"
	"github.com/joho/godotenv"
)

var (
	// ResourceBox is the packr box containing the resource file templates
	ResourceBox = packr.New("resource", "../../templates/resource")
	// ProjectBox is the packr box containing the project file templates
	ProjectBox = packr.New("project", "../../templates/project")
	// FunctionBox is the packr box containing the function file templates
	FunctionBox = packr.New("function", "../../templates/function")
	// SlsBox is the packr box containing the serverless.yml template
	SlsBox = packr.New("sls", "../../templates/sls")
)

// GetWorkingDir get the directory the current command is run out of
func GetWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return wd
}

// dirExists checks whether a folder with the given project name already exists
func dirExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// awsType returns the AWS datatype for a given golang type
func awsType(s string) string {
	switch strings.ToLower(s) {
	case "string", "time.Time", "*time.Time", "uuid.UUID":
		return "S"
	case "[]string":
		return "SS"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"byte", "rune",
		"float32", "float64", "complex64", "complex128":
		return "N"
	case "[]int", "[]int8", "[]int16", "[]int32", "[]int64",
		"[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64",
		"[]rune", "[]float32", "[]float64", "[]complex64", "[]complex128":
		return "NS"
	case "map[string]string", "map[string]int", "map[string]interface{}":
		return "M"
	case "[]byte":
		return "B"
	case "[][]byte":
		return "BS"
	case "bool":
		return "BOOL"

	default:
		return s
	}
}

// LoadTemplateFromBox loads a *text/template.Template from a packr.Box
func LoadTemplateFromBox(b *packr.Box, file string) *template.Template {
	// load string from template
	ts, err := b.FindString(file)
	if err != nil {
		log.Fatal(err)
	}

	// create new template with string
	t, err := template.New(file).Parse(ts)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

// appendIfMissing appends an element to a slice, if it doesn't contain the element already
func appendStringIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// RunCmd will run an OS command with the given arguments
func RunCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)

	err := execCmd(cmd)
	if err != nil {
		log.Fatalf("Executing %s failed with %s\n", name, err)
	}
}

// RunCmdWithEnv will run an OS command with the given arguments and an environment
func RunCmdWithEnv(envs []string, name string, args ...string) {
	cmdEnv := append(os.Environ(), envs...)
	cmd := exec.Command(name, args...)
	cmd.Env = cmdEnv

	err := execCmd(cmd)
	if err != nil {
		log.Fatalf("Executing %s failed with %s\n", name, err)
	}
}

func execCmd(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// func renderGopkg(config ResourceConfig) {
// 	log.Println("Generating Gopkg.toml")

// 	processed := map[string]bool{}

// 	// iterate resources
// 	for k, r := range config.Resources {
// 		path := filepath.Join(getPath(config, r), "Gopkg.toml")
// 		createGopkg(k, path)

// 		processed[k] = true
// 	}

// 	// generate serverless.yml for remaining functions
// 	for k, fs := range config.Functions {
// 		if !processed[k] {
// 			for _, f := range fs {
// 				path := filepath.Join(getPath(config, f), "Gopkg.toml")
// 				createGopkg(f.Name, path)
// 			}
// 		}
// 	}
// }

// func createGopkg(key, path string) {
// 	// load template
// 	t := loadTemplateFromBox(slsBox, "Gopkg.toml.tmpl")

// 	// create file only if it doesn't exist already
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		f, err := os.Create(path)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer f.Close()

// 		// execute template and save to file
// 		err = t.Execute(f, nil)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		log.Printf("Gopkg.toml generated for %s.\n", key)
// 	}
// }

func renderMakefile(config ResourceConfig) {
	log.Println("Generating Makefile...")

	// load Makefile template
	t := LoadTemplateFromBox(ProjectBox, "Makefile.tmpl")

	// open file and execute template
	f, err := os.OpenFile(filepath.Join(config.ProjectPath, "Makefile"), os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// execote template and save to file
	err = t.Execute(f, config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Makefile generated.")
}

// RenderSLS will render the serverless.yml file with a given ResourceConfig
func RenderSLS(config ResourceConfig) {
	log.Println("Generating serverless.yml...")

	processed := map[string]bool{}

	// generate serverless.yml for each resource
	for k := range config.Resources {
		path, resourceConfig := GetConfigForResource(k, config)
		GenerateSLS(filepath.Join(path, "serverless.yml"), resourceConfig)

		processed[k] = true
	}

	// generate serverless.yml for remaining functions
	for k, fs := range config.Functions {
		if !processed[k] {
			for _, f := range fs {
				path, functionConfig := GetConfigForFunction(k, f, config)
				GenerateSLS(filepath.Join(path, "serverless.yml"), functionConfig)
			}
		}
	}
}

// GetConfigForResource returns only the named resource ResourceConfig with the path information
func GetConfigForResource(k string, config ResourceConfig) (string, ResourceConfig) {
	r := config.Resources[k]
	path := getPath(config, r)

	// only handle current resource
	config.Resources = map[string]*Resource{
		k: r,
	}

	config.Functions = map[string][]*Function{
		k: config.Functions[k],
	}

	// load environment variables from .env
	config.Environments, _ = godotenv.Read(filepath.Join(config.ProjectPath, ".env"))

	return path, config
}

// GetConfigForFunction returns only the named function ResourceConfig with the path information
func GetConfigForFunction(k string, f *Function, config ResourceConfig) (string, ResourceConfig) {
	path := getPath(config, f)

	// only handle current function
	config.Resources = map[string]*Resource{}
	config.Functions = map[string][]*Function{
		"": []*Function{f},
	}

	// load environment variables from .env
	config.Environments, _ = godotenv.Read(filepath.Join(config.ProjectPath, ".env"))

	return path, config
}

func getPath(config ResourceConfig, i interface{}) string {
	path := filepath.Join(config.ProjectPath, "functions")

	switch t := i.(type) {
	case Function:
		path = filepath.Join(path, t.Name)
	case Resource:
		path = filepath.Join(path, t.Ident.String())
	}

	return path
}

// GenerateSLS generates the serverlss.yml for a given ResourceConfig
func GenerateSLS(path string, config ResourceConfig) {
	// load serverless.yml template
	t := LoadTemplateFromBox(SlsBox, "serverless.yml.tmpl")

	// create file
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// execute template and save to file
	err = t.Execute(f, config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("serverless.yml generated.")
}

func generateSAMTemplate(config ResourceConfig) {
	log.Println("Generating template.yml...")

	// load Makefile template
	t := LoadTemplateFromBox(ProjectBox, "template.yml.tmpl")

	// open file and execute template
	f, err := os.Create(filepath.Join(config.ProjectPath, "template.yml"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// execote template and save to file
	err = t.Execute(f, config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("template.yml generated.")
}

// RemoveFiles ...
func RemoveFiles(config ResourceConfig, resourceName string, function *flect.Ident) {
	// create the function folder
	folder := filepath.Join(config.ProjectPath, "functions", resourceName)
	if function != nil {
		folder = filepath.Join(folder, function.String())
	}

	err := os.RemoveAll(folder)
	if err != nil {
		log.Fatalf("Error deleting function folder %s: %s", function, err)
	}
}

// UpdateYMLs updates serverless.yml, Makefile, template.yml and create Gopkg.toml
func UpdateYMLs(config ResourceConfig, ignoreSLS bool) {
	// renderGopkg(config)
	renderMakefile(config)
	if !ignoreSLS {
		RenderSLS(config)
	}
	generateSAMTemplate(config)
}

func readDataFromFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetFuncName returns the generated function name for a given resource/ function group name and a functionName
func GetFuncName(resourceName, functionName string) string {
	ident := flect.New(resourceName)
	if ident.String() == "_" {
		return functionName
	}

	return functionName + "_" + ident.Singularize().String()

}
