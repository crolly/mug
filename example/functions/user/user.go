package user

import (
	"encoding/json"
	"fmt"
    "os"
	"strconv"

    "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

    "github.com/gofrs/uuid"
    )

// User defines the User model
type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`

	Address Address `json:"address"`
}

// Address defines the Address model
type Address struct {
	Street string `json:"street"`
	Zip string `json:"zip"`
	City string `json:"city"`
}



func connectDB() (*dynamodb.DynamoDB, string) {
	var sess *session.Session
	var tableName string

	local, err := strconv.ParseBool(os.Getenv("AWS_SAM_LOCAL"))
	if err != nil {
		local = false
	}
	// Create dynamo client object locally if running SAM CLI
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
		}))
		tableName = "users"
	} else {
		sess = session.Must(session.NewSession())
		tableName = os.Getenv("USER_TABLE_NAME")	
	}

	return dynamodb.New(sess), tableName
}

// Put extracts the User from JSON and writes it to DynamoDB
func Put(body string) (User, error) {
	svc, tableName := connectDB()

	// Marshall the requrest body
	var user User
	json.Unmarshal([]byte(body), &user)

	// Generate new UUID to store User in case user doesn't have one
	givenID, err := uuid.FromString(user.ID)
	if err != nil || givenID == uuid.Nil {
        id, _ := uuid.NewV4()
        user.ID = id.String()
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
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	return user, err
}

// Read gets the User from DynamoDB
func Read(id string) (User, error) {
    svc, tableName := connectDB()
	user := User{}

	// Perform the query
	fmt.Println("Trying to read from table: ", "users")
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
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

// Delete erases the User from DynamoDB
func Delete(id string) error {
    svc, tableName := connectDB()

	// Perform the delete
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// List returns the Users from DynamoDB
func List() ([]User, error){
    svc, tableName := connectDB()

    params := &dynamodb.ScanInput{
        TableName: aws.String(tableName),
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