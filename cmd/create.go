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

	//db represents the flag variable whether a db should be added
	// db bool
)

func init() {
	createCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	rootCmd.AddCommand(createCmd)

	// rootCmd.Flags().BoolVarP(&db, "database", "db", true, "Flag whether or not a db should be added as resource")
}

// createsProjectStructure creates the project structure with serverless.yml and mug.config.json
func createProjectStructure(projectName string) {
	// get working directory
	wd := getWorkingDir()

	// create root directory if it doesn't exist already, else fail
	if dirExists(filepath.Join(wd, projectName)) {
		log.Fatal("folder already exists")
	}
	os.Mkdir(projectName, 0755)

	// set data
	config := ResourceConfig{
		ProjectName: projectName,
	}

	// iterate over templates and execute
	for _, tmpl := range projectBox.List() {

		// create file
		fileName := fmt.Sprintf("%s/%s/%s", wd, projectName, strings.Replace(tmpl, ".tmpl", "", 1))
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
