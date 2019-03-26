// Copyright © 2019 Christian Rolly <mail@chromium-solutions.de>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// createCmd represents the create command
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Creates the boilerplate for your AWS Lambda for golang project.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createProjectStructure(args[0])
		},
	}

	region string
)

func init() {
	createCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	createCmd.Flags().StringVarP(&region, "region", "r", "eu-central-1", "Region the project will be deployed to (e.g. us-east-1 or eu-central-1)")

	rootCmd.AddCommand(createCmd)
}

// createsProjectStructure creates the project structure with serverless.yml and mug.config.json
func createProjectStructure(projectName string) {
	// create new config from project name
	config := newConfig(projectName)

	// create folder for project if it doesn't exist already
	if _, err := os.Stat(config.ProjectPath); !os.IsNotExist(err) {
		// projectPath exists already
		log.Fatal("folder already exists")
	}
	os.MkdirAll(config.ProjectPath, 0755)

	// iterate over templates and execute
	for _, tmpl := range projectBox.List() {

		// create file
		fileName := fmt.Sprintf("%s/%s", config.ProjectPath, strings.Replace(tmpl, ".tmpl", "", 1))
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// load template
		t := loadTemplateFromBox(projectBox, tmpl)

		// execute template and save to file
		err = t.Execute(f, config)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func newConfig(projectName string) ResourceConfig {
	pName, pPath, iPath := getPaths(projectName)

	config := ResourceConfig{
		ProjectName: pName,
		ProjectPath: pPath,
		ImportPath:  iPath,
		Region:      region,
	}

	return config
}

func getPaths(projectName string) (string, string, string) {
	projectPath, importPath := "", ""

	// environments GOPATH
	goPath := os.Getenv("GOPATH")
	if len(goPath) == 0 {
		log.Fatal("$GOPATH is not set")
	}
	srcPath := filepath.Join(goPath, "src")

	if strings.Contains(projectName, "/") {
		// project is created with full path to GOPATH src e.g. github.com/crolly/mug-example
		projectPath = filepath.Join(srcPath, projectName)
		importPath = projectName

		i := strings.LastIndex(projectName, "/")
		projectName = projectName[i+1 : len(projectName)]
	} else {
		// project is created with project name only
		wd := getWorkingDir()
		if filepathHasPrefix(wd, srcPath) {
			projectPath = filepath.Join(wd, projectName)
			importPath = strings.TrimPrefix(strings.ReplaceAll(projectPath, srcPath, ""), "/")
		} else {
			log.Fatal("You must either create the project inside of $GOPATH or provide the full path (e.g. github.com/crolly/mug-example")
		}
	}

	return projectName, projectPath, importPath
}

func filepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}
	return path[0:len(prefix)] == prefix

}
