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
	"encoding/json"
	"fmt"
	"io/ioutil"
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
			// update the yml files and Makefile with current config
			updateYMLs(readConfig(), noUpdate)
			// make debug binaries overwriting previous
			makeDebug()
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
	apiPort        string
)

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().BoolVarP(&remoteDebugger, "remoteDebugger", "r", false, "indicated whether you want to run a remote debugger")
	debugCmd.Flags().StringVarP(&debugPort, "debugPort", "p", "5986", "defines the remote port if remoteDebugger is true [default: 5986]")
	debugCmd.Flags().StringVarP(&apiPort, "apiPort", "a", "3000", "defines the port of local lambda api [default: 3000]")
	debugCmd.Flags().BoolVarP(&noUpdate, "disableYMLUpdate", "d", false, "Disable update of serverless.yml during execution")
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
			createTableForResource(svc, r)
		}

	}
}

func createTableForResource(svc *dynamodb.DynamoDB, resource Resource) {
	// get attributes
	attributes := []*dynamodb.AttributeDefinition{}
	for _, a := range resource.Attributes {
		attributes = append(attributes, &dynamodb.AttributeDefinition{
			AttributeName: aws.String(a.Ident.String()),
			AttributeType: aws.String(a.AwsType),
		})
	}

	// get keySchema
	keySchema := []*dynamodb.KeySchemaElement{}
	for t, k := range resource.KeySchema {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{
			AttributeName: aws.String(k),
			KeyType:       aws.String(t),
		})
	}

	// get throughput
	throughput := &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(10),
		WriteCapacityUnits: aws.Int64(10),
	}
	if len(resource.CapacityUnits) > 0 {
		throughput = &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(int64(resource.CapacityUnits["read"])),
			WriteCapacityUnits: aws.Int64(int64(resource.CapacityUnits["write"])),
		}
	}

	// set table name
	tableName := resource.Ident.Pluralize().String()

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
	args := []string{"local", "start-api", "-p", apiPort, "--docker-network", "lambda-local"}
	if remoteDebugger {
		ensureDebugger()
		args = append(args, "--debugger-path", "./dlv", "-d", debugPort, "--debug-args", "-delveAPI=2")
		log.Printf("Starting local API at port %s with debugger at %s...\n", apiPort, debugPort)
	}

	runCmd("sam", args...)
}

func ensureDebugger() {
	// build delve
	log.Println("Building dlv locally")
	env := []string{"GOARCH=amd64", "GOOS=linux"}
	runCmdWithEnv(env, "go", "build", "-o", "./dlv/dlv", "github.com/go-delve/delve/cmd/dlv")
}

//reads model definition for a resource
func getResourceForTable(table string) Resource {
	wd := getWorkingDir()

	configFile, err := os.Open(filepath.Join(wd, "mug.config.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config ResourceConfig
	json.Unmarshal(data, &config)

	return config.Resources[table]
}
