package main

import (
	"fmt"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/crolly/mug/example/functions/course"
)

// DeleteHandler handles the DELETE request and delete the course by given id
func DeleteHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model with params found in the path
	id := request.PathParameters["id"]
	fmt.Println("Path vars: ", id)
	err := course.Delete(id)
	if err != nil {
		panic(fmt.Sprintf("Failed to find course, %v", err))
	}

	msg := fmt.Sprintf("Deleted course with id: %s \n", id)
	return events.APIGatewayProxyResponse{Body: msg, StatusCode: 200}, nil
}

func main() {
	lambda.Start(DeleteHandler)
}