package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

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

// delete erases the Course from DynamoDB
func delete(id string) error {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Perform the delete
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(os.Getenv("COURSE_TABLE_NAME")),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// DeleteHandler handles the DELETE request and delete the course by given id
func DeleteHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model with params found in the path
	id := request.PathParameters["id"]
	fmt.Println("Path vars: ", id)
	err := delete(id)
	if err != nil {
		panic(fmt.Sprintf("Failed to find course, %v", err))
	}

	msg := fmt.Sprintf("Deleted course with id: %s \n", id)
	return events.APIGatewayProxyResponse{Body: msg, StatusCode: 200}, nil
}

func main() {
	lambda.Start(DeleteHandler)
}
