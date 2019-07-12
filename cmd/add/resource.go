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
	"strings"

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
			// TODO: check if resource exists already

			// instantiate new resource model and parse given attributes
			modelName := args[0]
			capacityUnits := map[string]int64{
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
			m := models.New(modelName, false, attributes, options)

			// get all imports
			m.Imports = m.GetImports()

			// add resource to mug.config.json
			mc, sc := m.GetConfigs()

			// render templates with data
			renderTemplates(mc, m)

			// write modelName.json, mug.config.json and serverless.yml for resource
			m.Write(mc.ProjectPath)
			mc.Write()
			sc.Write(mc.ProjectPath, modelName)
		},
	}

	attributes, keySchema, billingMode string
	noID, dates, softDelete            bool
	readUnits, writeUnits              int64
)

func init() {
	AddCmd.AddCommand(resourceCmd)
	resourceCmd.Flags().StringVarP(&attributes, "attributes", "a", "", "attributes of the resource")
	resourceCmd.Flags().BoolVarP(&noID, "noID", "n", false, "disable automatic generation of id attribute with type uuid")
	resourceCmd.Flags().BoolVarP(&dates, "addDates", "d", false, "automatically add createdAt and updatedAt attributes")
	resourceCmd.Flags().BoolVarP(&softDelete, "softDelete", "s", false, "automatically add deletedAt attribute")
	resourceCmd.Flags().StringVarP(&keySchema, "keySchema", "k", "id:HASH", "Key Schema definition for the DynamoDB Table Resource (only applied if noID flag is set to true")
	resourceCmd.Flags().StringVarP(&billingMode, "billingMode", "b", "provisioned", "Choose between 'provisioned' for ProvisionedThroughput (default) or 'ondemand'")
	resourceCmd.Flags().Int64VarP(&readUnits, "readUnits", "r", 1, "Set the ReadCapacityUnits if billingMode is set to ProvisionedThroughput")
	resourceCmd.Flags().Int64VarP(&writeUnits, "writeUnits", "w", 1, "Set the WriteCapacityUnits if billingMode is set to ProvisionedThroughput")

	// resourceCmd.Flags().BoolVarP(&noUpdate, "ignoreYMLUpdate", "i", false, "Ignore serverless.yml and template.yml during execution")

}

func renderTemplates(config models.MUGConfig, m models.Model) {
	// iterate over templates and execute
	for _, tmpl := range models.ResourceBox.List() {
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
		t := models.LoadTemplateFromBox(models.ResourceBox, tmpl)

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
