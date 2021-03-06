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

package remove

import (
	"os"
	"path/filepath"

	"github.com/crolly/mug/cmd/models"
	"github.com/spf13/cobra"
)

// RemoveCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove name",
	Short: "Remove resource or function group from your project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rName := args[0]

		mc := models.ReadMUGConfig()

		// delete resource from configuration
		mc.RemoveResource(rName)

		// remove the deployment if it exists
		if _, err := os.Stat(filepath.Join(mc.ProjectPath, "functions", rName, ".serverless")); !os.IsNotExist(err) {
			models.RunCmd("/bin/sh", "-c", "cd "+filepath.Join(mc.ProjectPath, "functions", rName)+"; sls remove")
		}

		// delete resource folder
		models.RemoveFiles(mc.ProjectPath, rName, "")
		mc.Write()
	},
}
