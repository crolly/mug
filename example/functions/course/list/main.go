package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/crolly/mug/example/functions/course"
)

// ListHandler handles the GET request and retrieves all courses from the database returning the items on success
func ListHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model
	courses, err := course.List()
	if err != nil {
		panic(fmt.Sprintf("Failed to find courses, %v", err))
	}

	// Log and return result
	jsonItems, _ := json.MarshalIndent(courses, "", "  ")
	stringItems := string(jsonItems)
	fmt.Println("Found items: ", stringItems)
	return events.APIGatewayProxyResponse{Body: stringItems, StatusCode: 200}, nil
}

func main() {
	lambda.Start(ListHandler)
}