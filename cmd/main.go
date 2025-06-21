package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleRequest is the function that AWS Lambda will invoke.
// It takes a request object from API Gateway and returns a response object.
func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request for path: %s", request.Path)

	responseBody := map[string]string{"status": "UP", "message": "Go serverless with AWS Lambda!"}

	// Marshal the response body into a JSON string.
	jsonBody, err := json.Marshal(responseBody)
	if err != nil {

		log.Printf("error marshaling response: %v", err)
		return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: http.StatusInternalServerError}, nil
	}

	// Return a successful response. API Gateway needs this specific struct.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
	}, nil
}

// main is the entrypoint for the Lambda function.
func main() {
	// The lambda.Start function is provided by the AWS SDK.
	// It takes our handler function and starts the Lambda execution environment.
	lambda.Start(handleRequest)
}
