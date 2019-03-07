package cmd

import (
	"log"
	"os"
	"strings"
	"text/template"

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
	case "string", "time.Time", "uuid.UUID":
		return "S"
	case "[]string":
		return "SS"
	case "int":
		return "N"
	case "[]int":
		return "NS"
	case "map[string]string", "map[string]int", "map[string]interface{}":
		return "M"
	case "[]byte":
		return "B"
	case "bool":
		return "BOOL"
	case "[][]byte":
		return "BS"
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
