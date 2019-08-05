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

package add

import (
	"log"
	"os"
	"path/filepath"

	"github.com/crolly/mug/cmd/models"

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
			capacityUnits := map[string]int64{
				"read":  readUnits,
				"write": writeUnits,
			}
			options := map[string]interface{}{
				"id":         generateID,
				"dates":      dates,
				"softDelete": softDelete,
				"keySchema":  keySchema,
				"billing":    billingMode,
				"capacity":   capacityUnits,
			}
			m := models.New(modelName, false, attributes, options)

			// get all imports
			m.Imports = m.GetImports()

			// add resource to mug.config.json
			mc, sc := m.GetConfigs()

			// check if resource exists already
			if _, err := os.Stat(filepath.Join(mc.ProjectPath, "functions", modelName)); !os.IsNotExist(err) {
				log.Fatalf("Function Group or Resource with the given name (%s) already exists. \n", modelName)
			}

			// render templates with data
			renderTemplates(mc, m)

			// write modelName.json, mug.config.json and serverless.yml for resource
			// m.Write(mc.ProjectPath)
			m.WriteOAS(mc, sc)
			mc.Write()
			sc.Write(mc.ProjectPath, modelName)
		},
	}

	attributes, keySchema, billingMode string
	generateID, dates, softDelete      bool
	readUnits, writeUnits              int64
)

func init() {
	AddCmd.AddCommand(resourceCmd)
	resourceCmd.Flags().StringVarP(&attributes, "attributes", "a", "", "attributes of the resource")
	resourceCmd.Flags().BoolVarP(&generateID, "generateID", "g", false, "automatic generation of id attribute with uuid")
	resourceCmd.Flags().BoolVarP(&dates, "addDates", "d", false, "automatically add createdAt and updatedAt attributes")
	resourceCmd.Flags().BoolVarP(&softDelete, "softDelete", "s", false, "automatically add deletedAt attribute")
	resourceCmd.Flags().StringVarP(&keySchema, "keySchema", "k", "id:HASH", "Key Schema definition for the DynamoDB Table Resource (not compatible with generateID)")
	resourceCmd.Flags().StringVarP(&billingMode, "billingMode", "b", "provisioned", "Choose between 'provisioned' for ProvisionedThroughput (default) or 'ondemand'")
	resourceCmd.Flags().Int64VarP(&readUnits, "readUnits", "r", 1, "Set the ReadCapacityUnits if billingMode is set to ProvisionedThroughput")
	resourceCmd.Flags().Int64VarP(&writeUnits, "writeUnits", "w", 1, "Set the WriteCapacityUnits if billingMode is set to ProvisionedThroughput")
}

func renderTemplates(config models.MUGConfig, m models.Model) {
	temps := []string{
		"create",
		"read",
		"update",
		"delete",
		"list",
		"model",
		"modelMocks",
	}

	data := map[string]interface{}{
		"Model":  m,
		"Config": config,
	}

	mName := m.Ident.Camelize().String()

	// iterate over resource templates and execute
	for _, t := range temps {
		// create the function folder for function templete (except model)
		folder := filepath.Join(config.ProjectPath, "functions", mName)
		if t == "model" {
			os.MkdirAll(folder, 0755)
			createFile(mName+".go", "model.tmpl", folder, data)
		} else if t == "modelMocks" {
			mockString := mName + "Mocks"
			folder = filepath.Join(config.ProjectPath, "mocks", mockString)
			os.MkdirAll(folder, 0755)
			createFile(mockString+".go", "modelMocks.tmpl", folder, data)
		} else {
			folder = filepath.Join(folder, t)
			os.MkdirAll(folder, 0755)
			for _, tf := range []string{"main", "main_test"} {
				createFile(tf+".go", filepath.Join(t, tf+".tmpl"), folder, data)
			}
		}
	}
}

func createFile(fName, tPath, folder string, data map[string]interface{}) {
	f, err := os.Create(filepath.Join(folder, fName))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// load template
	tmpl := models.LoadTemplateFromBox(models.ResourceBox, tPath)

	err = tmpl.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}
}
