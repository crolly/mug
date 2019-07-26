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

package debug

import (
	"log"

	"github.com/crolly/mug/cmd/models"

	"github.com/spf13/cobra"
)

var (
	// DebugCmd represents the debug command
	DebugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Start Local API for debugging",
		Long:  `This command generates a template.yml for aws-sam-cli and starts a local api to test or debug against`,
		Run: func(cmd *cobra.Command, args []string) {
			// get the config
			mc := models.ReadMUGConfig()

			list := models.GetList(mc.ProjectPath, debugList)

			// create lambda-local network if it doesn't exist already
			models.CreateLambdaNetwork()
			// start dynamodb-local
			models.StartLocalDynamoDB()
			// create tables for resources
			mc.CreateResourceTables(list, "debug")

			// render template.yml
			t := models.NewTemplate()
			for _, r := range list {
				sc := mc.ReadServerlessConfig(r)
				t.AddFunctionsFromServerlessConfig(sc, r)
			}
			t.Write(mc.ProjectPath)

			// make debug binaries overwriting previous
			mc.MakeDebug(list)

			// start aws-sam-cli local api
			startLocalAPI()
		},
	}

	remoteDebugger               bool
	debugPort, gwPort, debugList string
)

func init() {
	DebugCmd.Flags().BoolVarP(&remoteDebugger, "remoteDebugger", "r", false, "indicates whether you want to run a remote debugger (e.g. step through your code with VSCode)")
	DebugCmd.Flags().StringVarP(&debugPort, "debugPort", "d", "5986", "defines the remote port if remoteDebugger is true")
	DebugCmd.Flags().StringVarP(&gwPort, "gwPort", "g", "3000", "defines the port of local API Gateway")
	DebugCmd.Flags().StringVarP(&debugList, "list", "l", "all", "comma separated list of resources/ function groups to debug")
}

func startLocalAPI() {
	args := []string{"local", "start-api", "-p", gwPort, "--docker-network", "lambda-local"}
	if remoteDebugger {
		ensureDebugger()
		args = append(args, "--debugger-path", "./dlv", "-d", debugPort, "--debug-args", "-delveAPI=2")
		log.Printf("Starting local API at port %s with debugger at %s...\n", gwPort, debugPort)
	}

	models.RunCmd("sam", args...)
}

func ensureDebugger() {
	// build delve
	log.Println("Building dlv locally")
	env := []string{"GOARCH=amd64", "GOOS=linux"}
	models.RunCmdWithEnv(env, "go", "build", "-o", "./dlv/dlv", "github.com/go-delve/delve/cmd/dlv")
}
