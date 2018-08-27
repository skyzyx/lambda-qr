package main

import (
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
)

var statusCode int

// The API Gateway handler
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	spew.Config.Indent = "    "
	statusCode = int(200)
	cacheFrom := time.Now().Format(http.TimeFormat)

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":  "text/plain; charset=utf-8",
			"Last-Modified": cacheFrom,
			"Expires":       cacheFrom,
		},
		Body:       spew.Sdump(request),
		StatusCode: statusCode,
	}, nil
}

// The core function
func main() {
	lambda.Start(Handler)
}
