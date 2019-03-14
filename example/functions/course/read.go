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

// read gets the Course from DynamoDB
func read(id string) (Course, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	course := Course{}

	// Perform the query
	fmt.Println("Trying to read from table: ", "courses")
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("COURSE_TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return course, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &course)
	if err != nil {
		fmt.Println(err.Error())
		return course, err
	}

	return course, nil
}

// ReadHandler handles the GET request to retrieve a course from the database returning it on success
func ReadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model with params found in the path
	fmt.Println("Path vars: ", request.PathParameters["id"])
	course, err := read(request.PathParameters["id"])
	if err != nil {
		panic(fmt.Sprintf("Failed to find course, %v", err))
	}

	// Make sure the course isn't empty
	if course.ID == "" {
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
