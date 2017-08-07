package proxy

import (
	"encoding/json"
	"fmt"
)

var getEvent = `{
  "resource": "/{proxy+}",
  "path": "/pets/tobi",
  "httpMethod": "GET",
  "headers": {
    "Accept": "*/*",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "CA",
    "Host": "apex-ping.com",
    "User-Agent": "curl/7.48.0",
    "Via": "2.0 a44b4468444ef3ee67472bd5c5016098.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "VRxPGF8rOXD7xpRjAjseXfRrFD3wg-QPUHY6chzB9bR7pXlct1NTpg==",
    "X-Amzn-Trace-Id": "Root=1-59554c99-4375fc8705ccb554008b3aad",
    "X-Forwarded-For": "207.102.57.26, 54.182.214.69",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https"
  },
  "queryStringParameters": {
    "format": "json"
  },
  "pathParameters": {
    "proxy": "pets/tobi"
  },
  "stageVariables": {
    "env": "prod"
  },
  "requestContext": {
    "path": "/pets/tobi",
    "accountId": "111111111",
    "resourceId": "jcl9w3",
    "stage": "prod",
    "requestId": "344b184b-5cfc-11e7-8483-27dbb2d30a77",
    "identity": {
      "cognitoIdentityPoolId": null,
      "accountId": null,
      "cognitoIdentityId": null,
      "caller": null,
      "apiKey": "",
      "sourceIp": "207.102.57.26",
      "accessKey": null,
      "cognitoAuthenticationType": null,
      "cognitoAuthenticationProvider": null,
      "userArn": null,
      "userAgent": "curl/7.48.0",
      "user": null
    },
    "resourcePath": "/{proxy+}",
    "httpMethod": "GET",
    "apiId": "iwcgwgigca"
  },
  "body": null,
  "isBase64Encoded": false
}`

var postEvent = `{
  "resource": "/{proxy+}",
  "path": "/pets/tobi",
  "httpMethod": "POST",
  "headers": {
    "Accept": "*/*",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "CA",
    "content-type": "application/json",
    "Host": "apex-ping.com",
    "User-Agent": "curl/7.48.0",
    "Via": "2.0 b790a9f06b09414fec5d8b87e81d4b7f.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "_h1jFD3wjq6ZIyr8be6RS7Y7665jF9SjACmVodBMRefoQCs7KwTxMw==",
    "X-Amzn-Trace-Id": "Root=1-59554cc9-35de2f970f0fdf017f16927f",
    "X-Forwarded-For": "207.102.57.26, 54.182.214.86",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https"
  },
  "queryStringParameters": null,
  "pathParameters": {
    "proxy": "pets/tobi"
  },
  "requestContext": {
    "path": "/pets/tobi",
    "accountId": "111111111",
    "resourceId": "jcl9w3",
    "stage": "prod",
    "requestId": "50f6e0ce-5cfc-11e7-ada1-4f5cfe727f01",
    "identity": {
      "cognitoIdentityPoolId": null,
      "accountId": null,
      "cognitoIdentityId": null,
      "caller": null,
      "apiKey": "",
      "sourceIp": "207.102.57.26",
      "accessKey": null,
      "cognitoAuthenticationType": null,
      "cognitoAuthenticationProvider": null,
      "userArn": null,
      "userAgent": "curl/7.48.0",
      "user": null
    },
    "resourcePath": "/{proxy+}",
    "httpMethod": "POST",
    "apiId": "iwcgwgigca"
  },
  "body": "{ \"name\": \"Tobi\" }",
  "isBase64Encoded": false
}`

var postEventBinary = `{
  "resource": "/{proxy+}",
  "path": "/pets/tobi",
  "httpMethod": "POST",
  "headers": {
    "Accept": "*/*",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "CA",
    "content-type": "text/plain",
    "Host": "apex-ping.com",
    "User-Agent": "curl/7.48.0",
    "Via": "2.0 b790a9f06b09414fec5d8b87e81d4b7f.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "_h1jFD3wjq6ZIyr8be6RS7Y7665jF9SjACmVodBMRefoQCs7KwTxMw==",
    "X-Amzn-Trace-Id": "Root=1-59554cc9-35de2f970f0fdf017f16927f",
    "X-Forwarded-For": "207.102.57.26, 54.182.214.86",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https"
  },
  "queryStringParameters": null,
  "pathParameters": {
    "proxy": "pets/tobi"
  },
  "requestContext": {
    "path": "/pets/tobi",
    "accountId": "111111111",
    "resourceId": "jcl9w3",
    "stage": "prod",
    "requestId": "50f6e0ce-5cfc-11e7-ada1-4f5cfe727f01",
    "identity": {
      "cognitoIdentityPoolId": null,
      "accountId": null,
      "cognitoIdentityId": null,
      "caller": null,
      "apiKey": "",
      "sourceIp": "207.102.57.26",
      "accessKey": null,
      "cognitoAuthenticationType": null,
      "cognitoAuthenticationProvider": null,
      "userArn": null,
      "userAgent": "curl/7.48.0",
      "user": null
    },
    "resourcePath": "/{proxy+}",
    "httpMethod": "POST",
    "apiId": "iwcgwgigca"
  },
  "body": "SGVsbG8gV29ybGQ=",
  "isBase64Encoded": true
}`

func output(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s\n", string(b))
}

func ExampleInput_get() {
	var in Input
	json.Unmarshal([]byte(getEvent), &in)
	output(in)
	// Output:
	//   {
	//   "HTTPMethod": "GET",
	//   "Headers": {
	//     "Accept": "*/*",
	//     "CloudFront-Forwarded-Proto": "https",
	//     "CloudFront-Is-Desktop-Viewer": "true",
	//     "CloudFront-Is-Mobile-Viewer": "false",
	//     "CloudFront-Is-SmartTV-Viewer": "false",
	//     "CloudFront-Is-Tablet-Viewer": "false",
	//     "CloudFront-Viewer-Country": "CA",
	//     "Host": "apex-ping.com",
	//     "User-Agent": "curl/7.48.0",
	//     "Via": "2.0 a44b4468444ef3ee67472bd5c5016098.cloudfront.net (CloudFront)",
	//     "X-Amz-Cf-Id": "VRxPGF8rOXD7xpRjAjseXfRrFD3wg-QPUHY6chzB9bR7pXlct1NTpg==",
	//     "X-Amzn-Trace-Id": "Root=1-59554c99-4375fc8705ccb554008b3aad",
	//     "X-Forwarded-For": "207.102.57.26, 54.182.214.69",
	//     "X-Forwarded-Port": "443",
	//     "X-Forwarded-Proto": "https"
	//   },
	//   "Resource": "/{proxy+}",
	//   "PathParameters": {
	//     "proxy": "pets/tobi"
	//   },
	//   "Path": "/pets/tobi",
	//   "QueryStringParameters": {
	//     "format": "json"
	//   },
	//   "Body": "",
	//   "IsBase64Encoded": false,
	//   "StageVariables": {
	//     "env": "prod"
	//   },
	//   "RequestContext": {
	//     "APIID": "iwcgwgigca",
	//     "ResourceID": "jcl9w3",
	//     "RequestID": "344b184b-5cfc-11e7-8483-27dbb2d30a77",
	//     "HTTPMethod": "GET",
	//     "ResourcePath": "/{proxy+}",
	//     "AccountID": "111111111",
	//     "Stage": "prod",
	//     "Identity": {
	//       "APIKey": "",
	//       "AccountID": "",
	//       "UserAgent": "curl/7.48.0",
	//       "SourceIP": "207.102.57.26",
	//       "AccessKey": "",
	//       "Caller": "",
	//       "User": "",
	//       "UserARN": "",
	//       "CognitoIdentityID": "",
	//       "CognitoIdentityPoolID": "",
	//       "CognitoAuthenticationType": "",
	//       "CognitoAuthenticationProvider": ""
	//     }
	//   }
	// }
}

func ExampleInput_post() {
	var in Input
	json.Unmarshal([]byte(postEvent), &in)
	output(in)
	// Output:
	// {
	//   "HTTPMethod": "POST",
	//   "Headers": {
	//     "Accept": "*/*",
	//     "CloudFront-Forwarded-Proto": "https",
	//     "CloudFront-Is-Desktop-Viewer": "true",
	//     "CloudFront-Is-Mobile-Viewer": "false",
	//     "CloudFront-Is-SmartTV-Viewer": "false",
	//     "CloudFront-Is-Tablet-Viewer": "false",
	//     "CloudFront-Viewer-Country": "CA",
	//     "Host": "apex-ping.com",
	//     "User-Agent": "curl/7.48.0",
	//     "Via": "2.0 b790a9f06b09414fec5d8b87e81d4b7f.cloudfront.net (CloudFront)",
	//     "X-Amz-Cf-Id": "_h1jFD3wjq6ZIyr8be6RS7Y7665jF9SjACmVodBMRefoQCs7KwTxMw==",
	//     "X-Amzn-Trace-Id": "Root=1-59554cc9-35de2f970f0fdf017f16927f",
	//     "X-Forwarded-For": "207.102.57.26, 54.182.214.86",
	//     "X-Forwarded-Port": "443",
	//     "X-Forwarded-Proto": "https",
	//     "content-type": "application/json"
	//   },
	//   "Resource": "/{proxy+}",
	//   "PathParameters": {
	//     "proxy": "pets/tobi"
	//   },
	//   "Path": "/pets/tobi",
	//   "QueryStringParameters": null,
	//   "Body": "{ \"name\": \"Tobi\" }",
	//   "IsBase64Encoded": false,
	//   "StageVariables": null,
	//   "RequestContext": {
	//     "APIID": "iwcgwgigca",
	//     "ResourceID": "jcl9w3",
	//     "RequestID": "50f6e0ce-5cfc-11e7-ada1-4f5cfe727f01",
	//     "HTTPMethod": "POST",
	//     "ResourcePath": "/{proxy+}",
	//     "AccountID": "111111111",
	//     "Stage": "prod",
	//     "Identity": {
	//       "APIKey": "",
	//       "AccountID": "",
	//       "UserAgent": "curl/7.48.0",
	//       "SourceIP": "207.102.57.26",
	//       "AccessKey": "",
	//       "Caller": "",
	//       "User": "",
	//       "UserARN": "",
	//       "CognitoIdentityID": "",
	//       "CognitoIdentityPoolID": "",
	//       "CognitoAuthenticationType": "",
	//       "CognitoAuthenticationProvider": ""
	//     }
	//   }
	// }
}
