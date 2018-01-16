package proxy

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

// Start the proxy.
func Start(h http.Handler) {
	lambda.Start(func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		req, err := NewRequest(e)
		if err != nil {
			return nil, errors.Wrap(err, "creating proxy request")
		}

		res := NewResponse()

		h.ServeHTTP(res, req)

		// return events.APIGatewayProxyResponse{
		// 	Body:       "Hello from Native Go",
		// 	StatusCode: 200,
		// }, nil
		return res.End(), nil
	})
}
