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

package deploy

import (
	"path/filepath"

	"github.com/crolly/mug/cmd/models"
	"github.com/spf13/cobra"
)

var (
	// DeployCmd represents the deploy command
	DeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys the stack to AWS using serverless framework",
		Run: func(cmd *cobra.Command, args []string) {
			mc := models.ReadMUGConfig()

			list := models.GetList(mc.ProjectPath, buildList)

			if !noTest {
				// create lambda-local network if it doesn't exist already
				models.CreateLambdaNetwork()
				// start dynamodb-local
				models.StartLocalDynamoDB()
				// create tables for resources
				mc.CreateResourceTables(list, "test", force)
			}

			// build binaries
			mc.MakeBuild(list, !noTest)
			// deploy to AWS
			sls := "sls deploy --stage " + stage
			if len(profile) > 0 {
				sls = sls + " --aws-profile " + profile
			}
			for _, r := range list {
				models.RunCmd("/bin/sh", "-c", "cd "+filepath.Join(mc.ProjectPath, "functions", r)+";"+sls)
			}
		},
	}

	name, buildList, stage, profile string
	noUpdate, noTest, force         bool
)

func init() {
	DeployCmd.Flags().BoolVarP(&noUpdate, "ignoreYMLUpdate", "i", false, "Ignore update of serverless.yml during execution")
	DeployCmd.Flags().BoolVarP(&noTest, "do not test", "t", false, "Don't run tests before deploying")
	DeployCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the resource of function to deploy")
	DeployCmd.Flags().StringVarP(&buildList, "list", "l", "all", "comma separated list of resources/ function groups to debug")
	DeployCmd.Flags().StringVarP(&stage, "stage", "s", "dev", "define deployment stage")
	DeployCmd.Flags().StringVarP(&profile, "profile", "p", "", "define deployment profile")
	DeployCmd.Flags().BoolVarP(&force, "force overwrite", "f", false, "force overwrite existing tables (might be necessary if you changed you table definition - e.g. new index)")
}
