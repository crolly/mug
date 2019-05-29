package cmd

import (
	"encoding/json"
	"fmt"
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

// getWorkingDir get the directory the current command is run out of
func getWorkingDir() string {
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
func loadTemplateFromBox(b *packr.Box, file string) *template.Template {
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

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)

	err := execCmd(cmd)
	if err != nil {
		log.Fatalf("Executing %s failed with %s\n", name, err)
	}
}

func runCmdWithEnv(envs []string, name string, args ...string) {
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

func renderMakefile(config ResourceConfig) {
	log.Println("Generating Makefile...")

	// load Makefile template
	t := loadTemplateFromBox(projectBox, "Makefile.tmpl")

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

func writeResourceDefinition(m Model, config ResourceConfig) {
	json, _ := json.MarshalIndent(m, "", "  ")
	_ = ioutil.WriteFile(filepath.Join(config.ProjectPath, "functions", m.Name, fmt.Sprintf("%s.json", m.Name)), json, 0644)
}

func renderSLS(config ResourceConfig) {
	log.Println("Generating serverless.yml...")

	processed := map[string]bool{}

	// load Makefile template
	t := loadTemplateFromBox(slsBox, "serverless.yml.tmpl")

	// generate serverless.yml for each resource
	for k, r := range config.Resources {
		path, resourceConfig := getConfigForResource(k, t.Name(), r, config)
		generateSLS(path, resourceConfig, t)

		processed[k] = true
	}

	// generate serverless.yml for remaining functions
	for k, fs := range config.Functions {
		if !processed[k] {
			for _, f := range fs {
				path, functionConfig := getConfigForFunction(k, t.Name(), f, config)
				generateSLS(path, functionConfig, t)
			}
		}
	}
}

func getConfigForResource(k, fileName string, r Resource, config ResourceConfig) (string, ResourceConfig) {
	path := filepath.Join(config.ProjectPath, "functions", r.Ident.String(), strings.Replace(fileName, ".tmpl", "", 1))

	// only handle current resource
	config.Resources = map[string]Resource{
		k: r,
	}

	config.Functions = map[string][]Function{
		k: config.Functions[k],
	}

	// load environment variables from .env
	config.Environments, _ = godotenv.Read(filepath.Join(config.ProjectPath, ".env"))

	return path, config
}

func getConfigForFunction(k, fileName string, f Function, config ResourceConfig) (string, ResourceConfig) {
	path := filepath.Join(config.ProjectPath, "functions", f.Name, strings.Replace(fileName, ".tmpl", "", 1))

	// only handle current function
	config.Resources = map[string]Resource{}
	config.Functions = map[string][]Function{
		"": []Function{f},
	}

	// load environment variables from .env
	config.Environments, _ = godotenv.Read(filepath.Join(config.ProjectPath, ".env"))

	return path, config
}

func generateSLS(path string, config ResourceConfig, t *template.Template) {
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
	t := loadTemplateFromBox(projectBox, "template.yml.tmpl")

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

func removeFiles(config ResourceConfig, resourceName string, function *flect.Ident) {
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

// updateYMLs updates serverless.yml, Makefile, template.yml
func updateYMLs(config ResourceConfig) {
	renderMakefile(config)
	renderSLS(config)
	generateSAMTemplate(config)
}
