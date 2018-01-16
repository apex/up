package proxy

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Start the proxy.
func Start(h http.Handler) {
	lambda.Start(func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{
			Body:       "Hello from Native Go",
			StatusCode: 200,
		}, nil
	})
}
