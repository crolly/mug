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
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/spf13/cobra"
)

// debugCmd represents the debug command
var (
	debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "starts a aws-sam-cli local api",
		Long: `This command generates a template.yml for aws-sam-cli and starts
	a local api to test or debug against`,
		Run: func(cmd *cobra.Command, args []string) {
			// make debug binaries overwriting previous
			makeDebug()
			// generate new template.yml overwriting previous
			generateTemplate()
			// create lambda-local network if it doesn't exist already
			createLambdaNetwork()
			// start dynamodb-local
			startLocalDynamoDB()
			// create tables for resources
			createResourceTables()
			// start aws-sam-cli local api
			startLocalAPI()
		},
	}

	remoteDebugger bool
	debugPort      string
)

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().BoolVarP(&remoteDebugger, "remoteDebugger", "r", false, "indicated whether you want to run a remote debugger")
	debugCmd.Flags().StringVarP(&debugPort, "debugPort", "p", "5986", "defines the remote port if remoteDebugger is true [default: 5986]")
}

func makeDebug() {
	// check if Makefile exists in working directory
	wd := getWorkingDir()
	if _, err := os.Stat(filepath.Join(wd, "Makefile")); os.IsNotExist(err) {
		log.Fatal("no Makefile found - cannout build binaries")
	}

	// delete debug binaries if they exists
	if _, err := os.Stat(filepath.Join(wd, "debug")); os.IsNotExist(err) {
		runCmd("rm", "-rf", "debug/")
	}
	// run make debug
	log.Println("Building Debug Binaries...")
	runCmd("make", "debug")
}

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)

	err := execCmd(cmd)
	if err != nil {
		log.Fatalf("Executing %s failed with %s\n", name, err)
	}
}

func runCmdWithEnv(envs []string, name string, args ...string) {
	cmdEnv := append(os.Environ(), envs...)
	cmd := exec.Command(name, args...)
	cmd.Env = cmdEnv

	err := execCmd(cmd)
	if err != nil {
		log.Fatalf("Executing %s failed with %s\n", name, err)
	}
}

func execCmd(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func generateTemplate() {
	log.Println("Generating template.yml...")
	config := readConfig()

	// load Makefile template
	t := loadTemplateFromBox(projectBox, "template.yml.tmpl")

	// open file and execute template
	f, err := os.Create(filepath.Join(getWorkingDir(), "template.yml"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// execote template and save to file
	err = t.Execute(f, config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("template.yml generated.")
}

func createLambdaNetwork() {
	// check if network exists
	out, err := exec.Command("docker", "network", "ls", "--filter", "name=^lambda-local$", "--format", "{{.Name}}").Output()
	if err != nil {
		log.Fatal(err)
	}
	// create network if it doesn't exist
	if len(out) == 0 {
		log.Println("Creating lambda-local docker network")
		runCmd("docker", "network", "create", "lambda-local")
	} else {
		log.Println("Docker network lambda-local already exists, skipping creation...")
	}
}

func startLocalDynamoDB() {
	// check if container exists
	out, err := exec.Command("docker", "ps", "-a", "--filter", "network=lambda-local", "--filter", "ancestor=amazon/dynamodb-local", "--filter", "name=dynamodb", "--format", "{{.Status}}").Output()
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(string(out), "Exited") {
		log.Println("Restarting dynamodb-local container...")
		runCmd("docker", "restart", "dynamodb")
	}

	// create container if it doesn't exist already
	if len(out) == 0 {
		log.Println("Starting dynamodb-local...")
		wd := getWorkingDir()
		runCmd("docker", "run", "-v", fmt.Sprintf("%s:/dynamodb_local_db", wd), "-p", "8000:8000", "--net=lambda-local", "--name", "dynamodb", "-d", "amazon/dynamodb-local")
	}

	log.Println("dynamodb-local running.")
}

func createResourceTables() {
	// read the resource config
	config := readConfig()

	// create service to dynamodb
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String(config.Region),
	}))
	svc := dynamodb.New(sess)

	// iterate over resources
	for _, r := range config.Resources {
		// create the table for the resource
		input := &dynamodb.CreateTableInput{
			TableName: aws.String(r.Ident.Pluralize().ToLower().String()),
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: aws.String("B"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
		}

		out, err := svc.CreateTable(input)
		if err != nil {
			log.Fatalf("Error creating table %s: %s", r.Ident.Pluralize(), err)
		}

		log.Printf("Table %s created: %s", r.Ident.Pluralize(), out)
	}
}

func startLocalAPI() {
	args := []string{"local", "start-api", "--docker-network", "lambda-local"}
	if remoteDebugger {
		ensureDebugger()
		args = append(args, "--debugger-path", "./dlv", "-d", debugPort)
	}

	log.Printf("Starting local API at port %s...\n", debugPort)
	runCmd("sam", args...)
}

func ensureDebugger() {
	// build delve
	log.Println("Building dlv locally")
	env := []string{"GOARCH=amd64", "GOOS=linux"}
	runCmdWithEnv(env, "go", "build", "-o", "./dlv/dlv", "github.com/derekparker/delve/cmd/dlv")
}
