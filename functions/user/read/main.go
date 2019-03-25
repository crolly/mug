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

// read gets the User from DynamoDB
func read(id string) (User, error) {
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
	user := User{}

	// Perform the query
	fmt.Println("Trying to read from table: ", "users")
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("USER_TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				B: []byte(aws.StringValue(aws.String(id))),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	return user, nil
}

// ReadHandler handles the GET request to retrieve a user from the database returning it on success
func ReadHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model with params found in the path
	fmt.Println("Path vars: ", request.PathParameters["id"])
	user, err := read(request.PathParameters["id"])
	if err != nil {
		panic(fmt.Sprintf("Failed to find user, %v", err))
	}

	// Make sure the user isn't empty
	if user.ID == uuid.Nil {
		fmt.Println("Could not find user")
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, nil
	}

	// Log and return result
	jsonItem, _ := json.MarshalIndent(user, "", "  ")
	stringItem := string(jsonItem)
	fmt.Println("Found item: ", stringItem)
	return events.APIGatewayProxyResponse{Body: stringItem, StatusCode: 200}, nil
}

func main() {
	lambda.Start(ReadHandler)
}
