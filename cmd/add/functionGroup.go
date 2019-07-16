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

// functionGroupCmd represents the functionGroup command
var (
	functionGroupCmd = &cobra.Command{
		Use:   "functionGroup name [flags]",
		Short: "Adds a new function group, you can then add functions to with 'mug add function -r name [flags]'",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// instantiate new functionGroup
			groupName := args[0]

			// create new ServerlessConfig
			mc := models.ReadMUGConfig()

			// check if functionGroup exists already
			if _, err := os.Stat(filepath.Join(mc.ProjectPath, "functions", groupName)); !os.IsNotExist(err) {
				log.Fatalf("Function Group or Resource with the given name (%s) already exists. \n", groupName)
			}

			// create new function group
			sc := mc.NewServerlessConfig()

			// save to folder
			sc.Write(mc.ProjectPath, groupName)
		},
	}
)

func init() {
	AddCmd.AddCommand(functionGroupCmd)
}
