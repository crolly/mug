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
	"github.com/gobuffalo/flect"
	"github.com/spf13/cobra"
)

// rmfunctionCmd represents the rmfunction command
var rmfunctionCmd = &cobra.Command{
	Use:   "function functionName",
	Short: "Removes a function from a resource",
	Run: func(cmd *cobra.Command, args []string) {
		actual := args[0]
		function := flect.New(actual).Camelize()

		// get config and add function to it
		config := readConfig()
		config.RemoveFunction(resourceName, function.String())

		// remove files
		removeFiles(config, resourceName, &function)
		config.Write()

		// update serverless.yml, Makefile, template.yml
		renderMakefile()
		renderSLS()
		generateSAMTemplate()
	},
}

func init() {
	removeCmd.AddCommand(rmfunctionCmd)

	rmfunctionCmd.Flags().StringVarP(&resourceName, "resource", "r", "", "Name of the resource the function should be added to")

	rmfunctionCmd.MarkFlagRequired("resource")
}
