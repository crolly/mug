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
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

			// make debug binaries overwriting previous
			mc.MakeDebug(list)
			// create lambda-local network if it doesn't exist already
			createLambdaNetwork()
			// start dynamodb-local
			startLocalDynamoDB()
			// create tables for resources
			createResourceTables(mc, list)

			// render template.yml
			t := models.NewTemplate()
			for _, r := range list {
				sc := mc.ReadServerlessConfig(r)
				t.AddFunctionsFromServerlessConfig(sc, r)
			}
			t.Write(mc.ProjectPath)

			// start aws-sam-cli local api
			startLocalAPI()
		},
	}

	remoteDebugger               bool
	debugPort, gwPort, debugList string
)

func init() {
	DebugCmd.Flags().BoolVarP(&remoteDebugger, "remoteDebugger", "r", false, "indicates whether you want to run a remote debugger (e.g. step through your code with VSCode)")
	DebugCmd.Flags().StringVarP(&debugPort, "debugPort", "d", "5986", "defines the remote port if remoteDebugger is true [default: 5986]")
	DebugCmd.Flags().StringVarP(&gwPort, "gwPort", "g", "3000", "defines the port of local API Gateway [default: 3000]")
	DebugCmd.Flags().StringVarP(&debugList, "list", "l", "all", "comma separated list of resources/ function groups to debug [default: all]")
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
		models.RunCmd("docker", "network", "create", "lambda-local")
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
		models.RunCmd("docker", "restart", "dynamodb")
	}

	// create container if it doesn't exist already
	if len(out) == 0 {
		log.Println("Starting dynamodb-local...")
		wd := models.GetWorkingDir()
		models.RunCmd("docker", "run", "-v", fmt.Sprintf("%s:/dynamodb_local_db", wd), "-p", "8000:8000", "--net=lambda-local", "--name", "dynamodb", "-d", "amazon/dynamodb-local")
	}

	log.Println("dynamodb-local running.")
}

func createResourceTables(m models.MUGConfig, list []string) {
	// create service to dynamodb
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String(m.Region),
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
	for n, r := range m.Resources {
		if models.Contains(list, n) {
			sc := m.ReadServerlessConfig(n)
			rName := r.Ident.Pascalize().String() + "DynamoDbTable"
			props := sc.Resources.Resources[rName].Properties
			tableName := r.Ident.Pluralize().Camelize().String() + "-debug"

			if tables[tableName] {
				log.Printf("Table %s already exists, skipping creation...", tableName)
			} else {
				createTableForResource(svc, tableName, props)
			}
		}
	}
}

func createTableForResource(svc *dynamodb.DynamoDB, tableName string, props models.Properties) {
	// get attributes
	attributes := []*dynamodb.AttributeDefinition{}
	for _, a := range props.AttributeDefinitions {
		attributes = append(attributes, &dynamodb.AttributeDefinition{
			AttributeName: aws.String(a.AttributeName),
			AttributeType: aws.String(a.AttributeType),
		})
	}

	// get keySchema
	keySchema := []*dynamodb.KeySchemaElement{}
	for _, k := range props.KeySchema {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{
			AttributeName: aws.String(k.AttributeName),
			KeyType:       aws.String(k.KeyType),
		})
	}

	// get throughput
	throughput := &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(10),
		WriteCapacityUnits: aws.Int64(10),
	}
	if props.ProvisionedThroughput != nil {
		throughput = &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(props.ProvisionedThroughput.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(props.ProvisionedThroughput.WriteCapacityUnits),
		}
	}

	// create the table input for the resource
	input := &dynamodb.CreateTableInput{
		TableName:             aws.String(tableName),
		AttributeDefinitions:  attributes,
		KeySchema:             keySchema,
		ProvisionedThroughput: throughput,
	}

	out, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Error creating table %s: %s", tableName, err)
	}

	log.Printf("Table %s created: %s", tableName, out)
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
