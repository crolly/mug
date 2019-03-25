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

// list returns the Users from DynamoDB
func list() ([]User, error) {
	var sess *session.Session
	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		return []User{}, err
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

	params := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("USER_TABLE_NAME")),
	}
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var users []User
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return users, nil
}

// ListHandler handles the GET request and retrieves all users from the database returning the items on success
func ListHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Pass the call to the model
	users, err := list()
	if err != nil {
		panic(fmt.Sprintf("Failed to find users, %v", err))
	}

	// Make sure the users slice isn't empty
	if len(users) == 0 {
		fmt.Println("Could not find users")
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, nil
	}

	// Log and return result
	jsonItems, _ := json.MarshalIndent(users, "", "  ")
	stringItems := string(jsonItems)
	fmt.Println("Found items: ", stringItems)
	return events.APIGatewayProxyResponse{Body: stringItems, StatusCode: 200}, nil
}

func main() {
	lambda.Start(ListHandler)
}
