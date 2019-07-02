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
	"strings"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var (
	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Add authentication to resources and functions",
		Run: func(cmd *cobra.Command, args []string) {
			config := readConfig()
			// add user pool if provided to env
			if pool != "" {
				addPoolEnv(config, pool)
			}

			// add authentication to functions
			addAuth(config, excludes)

			// update serverless.yml with authentication information
			renderSLS(config)
		},
	}

	pool, excludes string
)

func init() {
	addCmd.AddCommand(authCmd)

	authCmd.Flags().StringVarP(&pool, "user pool", "p", "", "define the user pool to authenticate against")
	authCmd.Flags().StringVarP(&excludes, "excludes", "x", "", "list of functions or resources without authentication")

}

func addPoolEnv(config ResourceConfig, pool string) {
	f, err := os.OpenFile(filepath.Join(config.ProjectPath, ".env"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString("COGNITOUSERPOOL = " + pool + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func addAuth(config ResourceConfig, excludes string) {
	excludeSlice := strings.Split(excludes, ",")
	activeAuth := false

	for k, v := range config.Functions {
		if k != "" {
			// must be resource
			if !contains(excludeSlice, k) {
				for _, f := range v {
					f.Authentication = true
					activeAuth = true
				}
			}
		} else {
			// must be standalone functions --> iterate through function slice
			for _, f := range v {
				if !contains(excludeSlice, f.Name) {
					f.Authentication = true
					activeAuth = true
				}
			}
		}
	}

	config.Authentication = activeAuth

	// write back config
	config.Write()
}

func contains(s []string, v string) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}
