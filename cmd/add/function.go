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

	"github.com/gobuffalo/flect"

	"github.com/crolly/mug/cmd/models"

	"github.com/spf13/cobra"
)

// functionCmd represents the function command
var (
	functionCmd = &cobra.Command{
		Use:   "function functionName",
		Short: "Adds a function to a resource",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fName := args[0]

			// get config and add function to it
			mc := models.ReadMUGConfig()
			sc := mc.ReadServerlessConfig(rName)

			fn := &models.Function{
				Name:    models.GetFuncName(rName, fName),
				Path:    strings.TrimPrefix(path, "/"),
				Method:  strings.ToLower(method),
				Handler: fName,
			}

			sc.AddFunction(fn)
			sc.Write(mc.ProjectPath, rName)

			// generate files
			renderFunction(mc, rName, fName)
		},
	}

	rName  string
	path   string
	method string
)

func init() {
	AddCmd.AddCommand(functionCmd)

	functionCmd.Flags().StringVarP(&rName, "assign", "a", "generic", "Name of the resource or function group the function should be assigned to")
	functionCmd.Flags().StringVarP(&path, "path", "p", "", "Path the function will respond to e.g. /users")
	functionCmd.Flags().StringVarP(&method, "method", "m", "", "Method the function will respond to e.g. get")

	functionCmd.MarkFlagRequired("path")
	functionCmd.MarkFlagRequired("method")
}

func renderFunction(config models.MUGConfig, rName, fName string) {
	fIdent := flect.New(fName)

	// create the function folder
	folder := filepath.Join(config.ProjectPath, "functions", rName)
	funcFolder := filepath.Join(folder, fName)
	os.MkdirAll(funcFolder, 0755)

	// determine function templates and file names for resource or function group function
	funcNames := map[string]string{
		"blueprint_test.tmpl": "main_test.go",
	}
	resourceFunc := false
	if _, err := os.Stat(filepath.Join(folder, rName+".go")); os.IsNotExist(err) {
		funcNames["blueprint.tmpl"] = "main.go"
	} else {
		resourceFunc = true
		funcNames["resourceBlueprint.tmpl"] = "main.go"
	}

	data := map[string]interface{}{
		"ResourceName": rName,
		"Function":     fIdent,
		"Config":       config,
	}
	for tmpl, fn := range funcNames {
		// create file
		f, err := os.Create(filepath.Join(funcFolder, fn))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		t := models.LoadTemplateFromBox(models.FunctionBox, tmpl)

		// execute template and save to file
		err = t.Execute(f, data)
		if err != nil {
			log.Fatal(err)
		}
	}

	if resourceFunc {
		// also add function to resource file
		f, err := os.OpenFile(filepath.Join(folder, rName+".go"), os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		data := map[string]interface{}{
			"Function": fIdent,
		}
		t := models.LoadTemplateFromBox(models.FunctionBox, "resourceFunction.tmpl")
		err = t.Execute(f, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
