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
			generateSAMTemplate(readConfig())
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

	// run make debug
	log.Println("Building Debug Binaries...")
	runCmd("make", "debug")
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

	// get list of tables
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("Error during creation of resource tables: %s\n", err)
	}

	tables := make(map[string]bool)
	for _, t := range result.TableNames {
		tables[*t] = true
	}

	// iterate over resources
	for _, r := range config.Resources {
		tableName := r.Ident.Pluralize().ToLower().String()
		if tables[tableName] {
			log.Printf("Table %s already exists, skipping creation...", tableName)
		} else {
			createTableForResource(svc, tableName)
		}

	}
}

func createTableForResource(svc *dynamodb.DynamoDB, tableName string) {
	// create the table for the resource
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
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
		log.Fatalf("Error creating table %s: %s", tableName, err)
	}

	log.Printf("Table %s created: %s", tableName, out)
}

func startLocalAPI() {
	args := []string{"local", "start-api", "--docker-network", "lambda-local"}
	if remoteDebugger {
		ensureDebugger()
		args = append(args, "--debugger-path", "./dlv", "-d", debugPort)
		log.Printf("Starting local API at port 3000 with debugger at %s...\n", debugPort)
	}

	runCmd("sam", args...)
}

func ensureDebugger() {
	// build delve
	log.Println("Building dlv locally")
	env := []string{"GOARCH=amd64", "GOOS=linux"}
	runCmdWithEnv(env, "go", "build", "-o", "./dlv/dlv", "github.com/derekparker/delve/cmd/dlv")
}
