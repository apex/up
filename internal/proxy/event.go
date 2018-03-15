package proxy

// Identity is the identity information associated with the request.
type Identity struct {
	APIKey                        string
	AccountID                     string
	UserAgent                     string
	SourceIP                      string
	AccessKey                     string
	Caller                        string
	User                          string
	UserARN                       string
	CognitoIdentityID             string
	CognitoIdentityPoolID         string
	CognitoAuthenticationType     string
	CognitoAuthenticationProvider string
}

// RequestContext is the contextual information provided by API Gateway.
type RequestContext struct {
	APIID        string
	ResourceID   string
	RequestID    string
	HTTPMethod   string
	ResourcePath string
	AccountID    string
	Stage        string
	Identity     Identity
	Authorizer   map[string]string `json:"-"`
}

// Input is the input provided by API Gateway.
type Input struct {
	HTTPMethod            string
	Headers               map[string]string
	Resource              string
	PathParameters        map[string]string
	Path                  string
	QueryStringParameters map[string]string
	Body                  string
	IsBase64Encoded       bool
	StageVariables        map[string]string
	RequestContext        RequestContext
}

// Output is the output expected by API Gateway.
type Output struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers,omitempty"`
	Body            string            `json:"body,omitempty"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}
