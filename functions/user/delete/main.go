package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"time"

	"github.com/satori/go.uuid"
)

// User defines the User model
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active"`
	Email    string    `json:"email"`

	Address     Address       `json:"address"`
	Enrollments []Enrollments `json:"enrollments"`
}

// Address defines the Address model
type Address struct {
	Street string `json:"street"`
	Zip    string `json:"zip"`
	City   string `json:"city"`
}

// Enrollments defines the Enrollments model
type Enrollments struct {
	CourseID  uuid.UUID `json:"course_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// delete erases the User from DynamoDB
func delete(id string) error {
	var sess *session.Session
	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return err
	}
	// Create dynamo client object locally if running SAM CLI
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
		}))
	} else {
		sess = session.Must(session.NewSession())
	}
	svc := dynamodb.New(sess)

	// Perform the delete
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				B: []byte(aws.StringValue(aws.String(id))),
			},
		},
		TableName: aws.String(os.Getenv("USER_TABLE_NAME")),
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// DeleteHandler handles the DELETE request and delete the user by given id
func DeleteHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model with params found in the path
	id := request.PathParameters["id"]
	fmt.Println("Path vars: ", id)
	err := delete(id)
	if err != nil {
		panic(fmt.Sprintf("Failed to find user, %v", err))
	}

	msg := fmt.Sprintf("Deleted user with id: %s \n", id)
	return events.APIGatewayProxyResponse{Body: msg, StatusCode: 200}, nil
}

func main() {
	lambda.Start(DeleteHandler)
}
