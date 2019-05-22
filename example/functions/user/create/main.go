package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/crolly/mug/example/functions/user"
)

// CreateHandler handles the POST request and writes a user to the database returning the item on success
func CreateHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log body and pass to the model
	fmt.Println("Received body: ", request.Body)
	user, err := user.Put(request.Body)
	if err != nil {
		fmt.Println("Got error calling Put method")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	// Log and return result
	fmt.Println("Wrote item:  ", user)
	jsonItem, _ := json.MarshalIndent(user, "", "  ")
	stringItem := string(jsonItem)
	return events.APIGatewayProxyResponse{Body: stringItem, StatusCode: 200}, nil
}

func main() {
	lambda.Start(CreateHandler)
}