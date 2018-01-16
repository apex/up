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
		req, err := NewRequest(r)
		if err != nil {
			return events.APIGatewayProxyResponse{}, errors.Wrap(err, "creating proxy request")
		}

		res := NewResponse()
		h.ServeHTTP(res, req)
		return res.End(), nil
	})
}
