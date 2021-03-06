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
	"github.com/crolly/mug/cmd/models"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var (
	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Add authentication to a resource or function group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rName := args[0]
			mc := models.ReadMUGConfig()
			sc := mc.ReadServerlessConfig(rName)
			// add user pool if provided to env
			if pool != "" {
				sc.AddPoolEnv(mc, rName, pool)
			}

			// add authentication to functions
			sc.AddAuth(excludes)

			// update serverless.yml
			sc.Write(mc.ProjectPath, rName)
		},
	}

	pool, excludes string
)

func init() {
	AddCmd.AddCommand(authCmd)

	authCmd.Flags().StringVarP(&pool, "user pool", "p", "", "define the user pool to authenticate against")
	authCmd.Flags().StringVarP(&excludes, "excludes", "x", "", "list of functions in resource/ function group without authentication")

	cobra.MarkFlagRequired(authCmd.Flags(), "user pool")
}
