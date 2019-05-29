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
			capacityUnits := map[string]byte{
				"read":  readUnits,
				"write": writeUnits,
			}
			options := map[string]interface{}{
				"id":         !noID,
				"dates":      dates,
				"softDelete": softDelete,
				"keySchema":  keySchema,
				"billing":    billingMode,
				"capacity":   capacityUnits,
			}
			m := newModel(modelName, false, attributes, options)

			// get all imports
			m.Imports = m.getImports()

			// add resource to mug.config.json
			config := addResourceConfig(m)
			// render templates with data
			renderTemplates(config, m)

			// update the yml files and Makefile with current config
			updateYMLs(config)

			// write definition to resource folder
			writeResourceDefinition(m, config)
		},
	}

	attributes, keySchema, billingMode string
	noID, dates, softDelete            bool
	readUnits, writeUnits              byte
)

func init() {
	addCmd.AddCommand(resourceCmd)
	resourceCmd.Flags().StringVarP(&attributes, "attributes", "a", "", "attributes of the resource")
	resourceCmd.Flags().BoolVarP(&noID, "noID", "n", false, "disable automatic generation of id attribute with type uuid")
	resourceCmd.Flags().BoolVarP(&dates, "addDates", "d", false, "automatically add createdAt and updatedAt attributes")
	resourceCmd.Flags().BoolVarP(&softDelete, "softDelete", "s", false, "automatically add deletedAt attribute")
	resourceCmd.Flags().StringVarP(&keySchema, "keySchema", "k", "id:HASH", "Key Schema definition for the DynamoDB Table Resource (only applied if noID flag is set to true")
	resourceCmd.Flags().StringVarP(&billingMode, "billingMode", "b", "provisioned", "Choose between 'provisioned' for ProvisionedThroughput (default) or 'ondemand'")
	resourceCmd.Flags().Uint8VarP(&readUnits, "readUnits", "r", 1, "Set the ReadCapacityUnits if billingMode is set to ProvisionedThroughput")
	resourceCmd.Flags().Uint8VarP(&writeUnits, "writeUnits", "w", 1, "Set the WriteCapacityUnits if billingMode is set to ProvisionedThroughput")

}

func addResourceConfig(m Model) ResourceConfig {
	config := readConfig()

	singular := m.Ident.Singularize().String()
	plural := m.Ident.Pluralize().String()

	attributeDefinitions := map[string]AttributeDefinition{}
	for _, k := range m.KeySchema {
		a := m.Attributes[k]
		if len(a.Name) > 0 {
			attributeDefinitions[a.Name] = AttributeDefinition{
				Ident:   a.Ident,
				AwsType: a.AwsType,
			}
		}
	}

	resource := Resource{
		Ident:         flect.New(m.Name),
		Attributes:    attributeDefinitions,
		KeySchema:     m.KeySchema,
		CompositeKey:  m.CompositeKey,
		BillingMode:   m.BillingMode,
		CapacityUnits: m.CapacityUnits,
	}
	config.Resources[m.Name] = resource

	var path string
	if m.CompositeKey {
		path = fmt.Sprintf("%s/{%s}/{%s}", plural, m.KeySchema["HASH"], m.KeySchema["RANGE"])
	} else {
		path = fmt.Sprintf("%s/{%s}", plural, m.KeySchema["HASH"])
	}

	config.Functions[m.Name] = []Function{
		Function{Name: "create" + "_" + singular, Handler: "create", Path: plural, Method: "post"},
		Function{Name: "read" + "_" + singular, Handler: "read", Path: path, Method: "get"},
		Function{Name: "update" + "_" + singular, Handler: "update", Path: path, Method: "put"},
		Function{Name: "delete" + "_" + singular, Handler: "delete", Path: path, Method: "delete"},
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
