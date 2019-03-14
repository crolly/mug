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

// put extracts the Course from JSON and writes it to DynamoDB
func put(body string) (Course, error) {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Marshall the requrest body
	var course Course
	json.Unmarshal([]byte(body), &course)

	// Generate new UUID to store Course in case course doesn't have one
	if course.ID.String() == "" {
		id := uuid.Must(uuid.NewV4())
		course.ID = id
	}

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(course)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return course, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("COURSE_TABLE_NAME")),
	}
	_, err = svc.PutItem(input)
	return course, err
}

// UpdateHandler handles the PUT request and updates a course in the database returning the item on success
func UpdateHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log body and pass to the model
	fmt.Println("Received body: ", request.Body)
	course, err := put(request.Body)
	if err != nil {
		fmt.Println("Got error calling Put method")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	// Log and return result
	fmt.Println("Updated item: ", course)
	jsonItem, _ := json.MarshalIndent(course, "", "  ")
	stringItem := string(jsonItem)
	return events.APIGatewayProxyResponse{Body: stringItem, StatusCode: 200}, nil
}

func main() {
	lambda.Start(UpdateHandler)
}
