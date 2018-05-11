package proxy

// Identity is the identity information associated with the request.
type Identity struct {
	APIKey                        string `json:"apiKey"`
	AccountID                     string `json:"accountId"`
	UserAgent                     string `json:"userAgent"`
	SourceIP                      string `json:"sourceIp"`
	AccessKey                     string `json:"accessKey"`
	Caller                        string `json:"caller"`
	User                          string `json:"user"`
	UserARN                       string `json:"userARN"`
	CognitoIdentityID             string `json:"cognitoIdentityId"`
	CognitoIdentityPoolID         string `json:"cognitoIdentityPoolId"`
	CognitoAuthenticationType     string `json:"cognitoAuthenticationType"`
	CognitoAuthenticationProvider string `json:"cognitoAuthenticationProvider"`
}

// RequestContext is the contextual information provided by API Gateway.
type RequestContext struct {
	APIID        string            `json:"apiId"`
	ResourceID   string            `json:"resourceId"`
	RequestID    string            `json:"requestId"`
	HTTPMethod   string            `json:"-"`
	ResourcePath string            `json:"-"`
	AccountID    string            `json:"accountId"`
	Stage        string            `json:"stage"`
	Identity     Identity          `json:"identity"`
	Authorizer   map[string]string `json:"authorizer"`
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
