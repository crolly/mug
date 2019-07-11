// Copyright © 2000 Christian Rolly <mail@chromium-solutions.de>
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

package add

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/crolly/mug/cmd/models"
	"github.com/gobuffalo/flect"

	"github.com/spf13/cobra"
)

// functionCmd represents the function command
var (
	functionCmd = &cobra.Command{
		Use:   "function functionName",
		Short: "Adds a function to a resource",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			actual := args[0]
			function := flect.New(actual).Camelize()

			// get config and add function to it
			config := models.ReadConfig()
			resourceName, f := config.AddFunction(resourceName, function.String(), strings.TrimPrefix(path, "/"), strings.ToLower(method))

			// generate files
			renderFunction(config, resourceName, function)
			config.Write()

			// update the yml files and Makefile with current config
			// updateYMLs(readConfig(), noUpdate)
			if noUpdate {
				path, functionConfig := models.GetConfigForFunction("", f, config)
				models.GenerateSLS(path, functionConfig)
			} else {
				models.UpdateYMLs(config, noUpdate)
			}
		},
	}

	resourceName string
	path         string
	method       string
)

func init() {
	AddCmd.AddCommand(functionCmd)

	functionCmd.Flags().StringVarP(&resourceName, "resource", "r", "", "Name of the resource the function should be added to")
	functionCmd.Flags().StringVarP(&path, "path", "p", "", "Path the function will respond to e.g. /users")
	functionCmd.Flags().StringVarP(&method, "method", "m", "", "Method the function will respond to e.g. get")

	functionCmd.Flags().BoolVarP(&noUpdate, "ignoreYMLUpdate", "i", false, "Ignore update of serverless.yml and template.yml during execution")

	functionCmd.MarkFlagRequired("path")
	functionCmd.MarkFlagRequired("method")
}

func renderFunction(config models.ResourceConfig, resourceName string, function flect.Ident) {

	// create the function folder
	folder := filepath.Join(config.ProjectPath, "functions", resourceName)
	funcFolder := filepath.Join(folder, function.String())
	os.MkdirAll(funcFolder, 0755)

	// create main.go file
	f, err := os.Create(filepath.Join(funcFolder, "main.go"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// get blueprint template
	var t *template.Template
	resourceFunc := false
	if resourceName != "" {
		t = models.LoadTemplateFromBox(models.FunctionBox, "resourceBlueprint.tmpl")
		resourceFunc = true
	} else {
		t = models.LoadTemplateFromBox(models.FunctionBox, "blueprint.tmpl")
	}

	// execute template and save to file
	data := map[string]interface{}{
		"ResourceName": resourceName,
		"Function":     function,
		"Config":       config,
	}
	err = t.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}

	if resourceFunc {
		// also add function to resource file
		f, err := os.OpenFile(filepath.Join(folder, resourceName+".go"), os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		data := map[string]interface{}{
			"Function": function,
		}
		t := models.LoadTemplateFromBox(models.FunctionBox, "resourceFunction.tmpl")
		err = t.Execute(f, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
