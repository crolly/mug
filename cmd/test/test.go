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

package test

import (
	"fmt"

	"github.com/crolly/mug/cmd/models"

	"github.com/spf13/cobra"
)

var (
	// TestCmd represents the debug command
	TestCmd = &cobra.Command{
		Use:   "test",
		Short: "Run go tests for your project",
		Run: func(cmd *cobra.Command, args []string) {
			// get the config
			mc := models.ReadMUGConfig()

			list := models.GetList(mc.ProjectPath, list)

			// create lambda-local network if it doesn't exist already
			models.CreateLambdaNetwork()
			// start dynamodb-local
			models.StartLocalDynamoDB()
			// create tables for resources
			mc.CreateResourceTables(list, "test", force)

			env := []string{"MODE=test"}
			for _, r := range list {
				t := "go test ./functions/" + r + "/... -cover"
				fmt.Println(t)
				models.RunCmdWithEnv(env, "/bin/sh", "-c", t)
			}
		},
	}

	list  string
	force bool
)

func init() {
	TestCmd.Flags().StringVarP(&list, "list", "l", "all", "comma separated list of resources/ function groups to debug")
	TestCmd.Flags().BoolVarP(&force, "force overwrite", "f", false, "force overwrite existing tables (might be necessary if you changed you table definition - e.g. new index)")
}
