package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofrs/uuid"

	"github.com/crolly/mug/example/functions/course"
)

// ReadHandler handles the GET request to retrieve a course from the database returning it on success
func ReadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model with params found in the path
	fmt.Println("Path vars: ", request.PathParameters["id"])
	course, err := course.Read(request.PathParameters["id"])
	if err != nil {
		panic(fmt.Sprintf("Failed to find course, %v", err))
	}

	// Make sure the course isn't empty
	if uuid.Must(uuid.FromString(course.ID)) == uuid.Nil {
		fmt.Println("Could not find course")
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, nil
	}

	// Log and return result
	jsonItem, _ := json.MarshalIndent(course, "", "  ")
	stringItem := string(jsonItem)
	fmt.Println("Found item: ", stringItem)
	return events.APIGatewayProxyResponse{Body: stringItem, StatusCode: 200}, nil
}

func main() {
	lambda.Start(ReadHandler)
}