package main

import (
	"fmt"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/crolly/mug/example/functions/user"
)

// DeleteHandler handles the DELETE request and delete the user by given id
func DeleteHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model with params found in the path
	id := request.PathParameters["id"]
	fmt.Println("Path vars: ", id)
	err := user.Delete(id)
	if err != nil {
		panic(fmt.Sprintf("Failed to find user, %v", err))
	}

	msg := fmt.Sprintf("Deleted user with id: %s \n", id)
	return events.APIGatewayProxyResponse{Body: msg, StatusCode: 200}, nil
}

func main() {
	lambda.Start(DeleteHandler)
}