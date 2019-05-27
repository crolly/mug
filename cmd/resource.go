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

	"github.com/gobuffalo/flect"

	"github.com/spf13/cobra"
)

// resourceCmd represents the resource command
var (
	resourceCmd = &cobra.Command{
		Use:   "resource name [flags]",
		Short: "Adds CRUDL functions for the defined resource",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// instantiate new resource model and parse given attributes
			modelName := args[0]
			m := newModel(modelName, false, attributes, !noID, addDates)

			// get all imports
			m.Imports = m.getImports()

			// add resource to mug.config.json
			config := addResourceConfig(m)
			// render templates with data
			renderTemplates(config, m)

			// update the yml files and Makefile with current config
			updateYMLs(config)

			// write definition to resource folder
			writeResourceDefinition(m, modelName)
		},
	}

	attributes string
	noID       bool
	addDates   bool
)

func init() {
	addCmd.AddCommand(resourceCmd)
	resourceCmd.Flags().StringVarP(&attributes, "attributes", "a", "", "attributes of the resource")
	resourceCmd.Flags().BoolVarP(&noID, "noID", "n", false, "automatically generate id attribute as hash key for resource")
	resourceCmd.Flags().BoolVarP(&addDates, "addDates", "d", false, "automatically add createdAt and updatedAt attributes")
}

func addResourceConfig(m Model) ResourceConfig {
	config := readConfig()

	singular := m.Ident.Singularize().String()
	plural := m.Ident.Pluralize().String()

	resource := Resource{
		Ident: flect.New(m.Name),
	}
	config.Resources[m.Name] = resource

	config.Functions[m.Name] = []Function{
		Function{Name: "create" + "_" + singular, Handler: "create", Path: plural, Method: "post"},
		Function{Name: "read" + "_" + singular, Handler: "read", Path: fmt.Sprintf("%s/{id}", plural), Method: "get"},
		Function{Name: "update" + "_" + singular, Handler: "update", Path: fmt.Sprintf("%s/{id}", plural), Method: "put"},
		Function{Name: "delete" + "_" + singular, Handler: "delete", Path: fmt.Sprintf("%s/{id}", plural), Method: "delete"},
		Function{Name: "list" + "_" + plural, Handler: "list", Path: plural, Method: "get"},
	}

	config.Write()

	return config
}

func renderTemplates(config ResourceConfig, m Model) {
	// iterate over templates and execute
	for _, tmpl := range resourceBox.List() {
		// create the function folder for function templete (except model)
		folder := filepath.Join(config.ProjectPath, "functions", m.Ident.Camelize().String())
		if tmpl != "model.go.tmpl" {
			folder = filepath.Join(folder, strings.Replace(tmpl, ".tmpl", "", 1))
		}
		os.MkdirAll(folder, 0755)

		// create files
		file := "main.go"
		if tmpl == "model.go.tmpl" {
			file = m.Ident.Camelize().String() + ".go"
		}
		f, err := os.Create(filepath.Join(folder, file))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// load template
		t := loadTemplateFromBox(resourceBox, tmpl)

		// execute template and save to file
		data := map[string]interface{}{
			"Model":  m,
			"Config": config,
		}
		err = t.Execute(f, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
