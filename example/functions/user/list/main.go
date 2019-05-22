package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/crolly/mug/example/functions/user"
)

// ListHandler handles the GET request and retrieves all users from the database returning the items on success
func ListHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model
	users, err := user.List()
	if err != nil {
		panic(fmt.Sprintf("Failed to find users, %v", err))
	}

	// Log and return result
	jsonItems, _ := json.MarshalIndent(users, "", "  ")
	stringItems := string(jsonItems)
	fmt.Println("Found items: ", stringItems)
	return events.APIGatewayProxyResponse{Body: stringItems, StatusCode: 200}, nil
}

func main() {
	lambda.Start(ListHandler)
}