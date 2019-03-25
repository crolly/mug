package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

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

// put extracts the User from JSON and writes it to DynamoDB
func put(body string) (User, error) {
	var sess *session.Session
	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return User{}, err
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

	// Marshall the requrest body
	var user User
	json.Unmarshal([]byte(body), &user)

	// Generate new UUID to store User in case user doesn't have one
	if user.ID == uuid.Nil {
		id := uuid.NewV4()
		user.ID = id
	}

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println("Got error marshalling map:", err.Error())
		return user, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("USER_TABLE_NAME")),
	}
	_, err = svc.PutItem(input)
	return user, err
}

// CreateHandler handles the POST request and writes a user to the database returning the item on success
func CreateHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log body and pass to the model
	fmt.Println("Received body: ", request.Body)
	user, err := put(request.Body)
	if err != nil {
		fmt.Println("Got error calling Put method")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	// Log and return result
	fmt.Println("Wrote item:  ", user)
	jsonItem, _ := json.MarshalIndent(user, "", "  ")
	stringItem := string(jsonItem)
	return events.APIGatewayProxyResponse{Body: stringItem, StatusCode: 200}, nil
}

func main() {
	lambda.Start(CreateHandler)
}
