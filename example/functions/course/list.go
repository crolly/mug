package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/satori/go.uuid"
)

// Course defines the Course model
type Course struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Subtitle    string    `json:"subtitle"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}

// list returns the Courses from DynamoDB
func list() ([]Course, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("COURSE_TABLE_NAME")),
	}
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var courses []Course
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &courses)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return courses, nil
}

// ListHandler handles the GET request and retrieves all courses from the database returning the items on success
func ListHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model
	courses, err := list()
	if err != nil {
		panic(fmt.Sprintf("Failed to find courses, %v", err))
	}

	// Make sure the courses slice isn't empty
	if len(courses) == 0 {
		fmt.Println("Could not find courses")
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, nil
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
