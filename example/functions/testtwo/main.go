package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Testtwo function description
func TesttwoHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  	// Log and return result
	jsonItem, err := json.MarshalIndent(map[string]string{"msg": "Testtwo invoked successfully"}, "", "  ")
    if err != nil {
        fmt.Println("Error occured")
        return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil 
    }
	stringItem := string(jsonItem)
	return events.APIGatewayProxyResponse{Body: stringItem, StatusCode: 200}, nil
}

func main() {
	lambda.Start(TesttwoHandler)
}