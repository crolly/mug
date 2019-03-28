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

func renderMakefile() {
	log.Println("Generating Makefile...")
	config := readConfig()

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

func writeResourceDefinition(m Model, name string) {
	wd := getWorkingDir()

	json, _ := json.MarshalIndent(m, "", "  ")
	_ = ioutil.WriteFile(filepath.Join(wd, "functions", name, fmt.Sprintf("%s.json", name)), json, 0644)
}

func renderSLS() {
	log.Println("Generating serverless.yml...")
	config := readConfig()

	// load Makefile template
	t := loadTemplateFromBox(projectBox, "serverless.yml.tmpl")

	// open file and execute template
	f, err := os.OpenFile(filepath.Join(config.ProjectPath, strings.Replace(t.Name(), ".tmpl", "", 1)), os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// execote template and save to file
	err = t.Execute(f, config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("serverless.yml generated.")
}

func generateSAMTemplate() {
	log.Println("Generating template.yml...")
	config := readConfig()

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
