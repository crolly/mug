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
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys the stack to AWS using serverless framework",
	Run: func(cmd *cobra.Command, args []string) {
		// update the yml files and Makefile with current config
		updateYMLs(readConfig())
		// build binaries
		makeBuild()
		// deploy to AWS
		deploy()
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func makeBuild() {
	// check if Makefile exists in working directory
	wd := getWorkingDir()
	if _, err := os.Stat(filepath.Join(wd, "Makefile")); os.IsNotExist(err) {
		log.Fatal("no Makefile found - cannout build binaries")
	}

	// run make debug
	log.Println("Building Binaries...")
	runCmd("make", "build")
}

func deploy() {
	// check if Makefile exists in working directory
	wd := getWorkingDir()
	if _, err := os.Stat(filepath.Join(wd, "serverless.yml")); os.IsNotExist(err) {
		log.Fatal("no serverless.yml found - cannout build binaries")
	}

	// run make debug
	log.Println("Deploying ...")
	runCmd("sls", "deploy")
}
