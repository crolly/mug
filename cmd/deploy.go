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
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var (
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys the stack to AWS using serverless framework",
		Run: func(cmd *cobra.Command, args []string) {
			config := readConfig()
			if disableUpdate {
				// render only Makefile
				renderMakefile(config)
			} else {
				// update the yml files and Makefile with current config
				updateYMLs(config, noUpdate)
			}

			// build binaries
			makeBuild(config)
			// deploy to AWS
			deploy(config)
		},
	}

	disableUpdate bool
	name          string
)

func init() {
	deployCmd.Flags().BoolVarP(&noUpdate, "disableYMLUpdate", "d", false, "Disable update of serverless.yml during execution")
	deployCmd.Flags().BoolVarP(&disableUpdate, "disableYMLUpdate", "d", false, "Disable update of template.yml and serverless.yml during execution")
	deployCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the resource of function to deploy.")
	rootCmd.AddCommand(deployCmd)
}

func makeBuild(config ResourceConfig) {
	// check if Makefile exists in working directory
	if _, err := os.Stat(filepath.Join(config.ProjectPath, "Makefile")); os.IsNotExist(err) {
		log.Fatal("no Makefile found - cannout build binaries")
	}

	// run make debug
	log.Println("Building Binaries...")
	runCmd("make", "build")
}

func deploy(config ResourceConfig) {
	dir := filepath.Join(config.ProjectPath, "functions")
	// check if only single resource of function should be deployed
	if name != "" {
		// deploy single resource or functions
		deployResourceOrFunction(dir, name)
	} else {
		// deploy all
		for k, fs := range config.Functions {
			if k == "" {
				for _, f := range fs {
					deployResourceOrFunction(dir, f.Name)
				}
			} else {
				deployResourceOrFunction(dir, k)
			}
		}
	}
}

func deployResourceOrFunction(dir, name string) {
	if _, err := os.Stat(filepath.Join(dir, name, "serverless.yml")); os.IsNotExist(err) {
		log.Fatal("no serverless.yml found - cannout deploy")
	}

	// run make debug
	log.Println("Deploying ...")
	runCmd("sls", "deploy")
}
