package resources

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/apex/up"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/tj/assert"
)

// keys returns keys from a map.
func keys(m Map) (v []string) {
	for k := range m {
		v = append(v, k)
	}
	return
}

// getResource returns resource by name.
func getResource(c *Config, name string) Map {
	tmpl := New(c)
	r := tmpl["Resources"].(Map)

	v, ok := r[name].(Map)
	if !ok {
		k := strings.Join(keys(r), "\n  - ")
		panic(fmt.Sprintf("resource %q does not exist in:\n\n  - %s", name, k))
	}

	return v
}

// dump a resource to stdout.
func dump(c *Config, name string) {
	r := getResource(c, name)
	{
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(r)
	}
}

func Example_api() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "Api")
	// Output:
	// {
	//   "Properties": {
	//     "BinaryMediaTypes": [
	//       "*/*"
	//     ],
	//     "Description": "Managed by Up.",
	//     "Name": {
	//       "Ref": "Name"
	//     }
	//   },
	//   "Type": "AWS::ApiGateway::RestApi"
	// }
}

func Example_apiRootMethod() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "ApiRootMethod")
	// Output:
	//  {
	//   "Properties": {
	//     "AuthorizationType": "NONE",
	//     "HttpMethod": "ANY",
	//     "Integration": {
	//       "IntegrationHttpMethod": "POST",
	//       "Type": "AWS_PROXY",
	//       "Uri": {
	//         "Fn::Join": [
	//           "",
	//           [
	//             "arn:aws:apigateway:",
	//             {
	//               "Ref": "AWS::Region"
	//             },
	//             ":lambda:path/2015-03-31/functions/",
	//             {
	//               "Fn::Join": [
	//                 ":",
	//                 [
	//                   "arn",
	//                   "aws",
	//                   "lambda",
	//                   {
	//                     "Ref": "AWS::Region"
	//                   },
	//                   {
	//                     "Ref": "AWS::AccountId"
	//                   },
	//                   "function",
	//                   {
	//                     "Fn::Join": [
	//                       ":",
	//                       [
	//                         {
	//                           "Ref": "FunctionName"
	//                         },
	//                         "${stageVariables.qualifier}"
	//                       ]
	//                     ]
	//                   }
	//                 ]
	//               ]
	//             },
	//             "/invocations"
	//           ]
	//         ]
	//       }
	//     },
	//     "ResourceId": {
	//       "Fn::GetAtt": [
	//         "Api",
	//         "RootResourceId"
	//       ]
	//     },
	//     "RestApiId": {
	//       "Ref": "Api"
	//     }
	//   },
	//   "Type": "AWS::ApiGateway::Method"
	// }
}

func Example_apiProxyResource() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "ApiProxyResource")
	// Output:
	// {
	//   "Properties": {
	//     "ParentId": {
	//       "Fn::GetAtt": [
	//         "Api",
	//         "RootResourceId"
	//       ]
	//     },
	//     "PathPart": "{proxy+}",
	//     "RestApiId": {
	//       "Ref": "Api"
	//     }
	//   },
	//   "Type": "AWS::ApiGateway::Resource"
	// }
}

func Example_apiProxyMethod() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "ApiProxyMethod")
	// Output:
	// {
	//   "Properties": {
	//     "AuthorizationType": "NONE",
	//     "HttpMethod": "ANY",
	//     "Integration": {
	//       "IntegrationHttpMethod": "POST",
	//       "Type": "AWS_PROXY",
	//       "Uri": {
	//         "Fn::Join": [
	//           "",
	//           [
	//             "arn:aws:apigateway:",
	//             {
	//               "Ref": "AWS::Region"
	//             },
	//             ":lambda:path/2015-03-31/functions/",
	//             {
	//               "Fn::Join": [
	//                 ":",
	//                 [
	//                   "arn",
	//                   "aws",
	//                   "lambda",
	//                   {
	//                     "Ref": "AWS::Region"
	//                   },
	//                   {
	//                     "Ref": "AWS::AccountId"
	//                   },
	//                   "function",
	//                   {
	//                     "Fn::Join": [
	//                       ":",
	//                       [
	//                         {
	//                           "Ref": "FunctionName"
	//                         },
	//                         "${stageVariables.qualifier}"
	//                       ]
	//                     ]
	//                   }
	//                 ]
	//               ]
	//             },
	//             "/invocations"
	//           ]
	//         ]
	//       }
	//     },
	//     "ResourceId": {
	//       "Ref": "ApiProxyResource"
	//     },
	//     "RestApiId": {
	//       "Ref": "Api"
	//     }
	//   },
	//   "Type": "AWS::ApiGateway::Method"
	// }
}

func Example_apiDeploymentDevelopment() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "ApiDeploymentDevelopment")
	// Output:
	// {
	//   "DependsOn": [
	//     "ApiRootMethod",
	//     "ApiProxyMethod",
	//     "ApiFunctionAliasDevelopment"
	//   ],
	//   "Properties": {
	//     "RestApiId": {
	//       "Ref": "Api"
	//     },
	//     "StageDescription": {
	//       "Variables": {
	//         "qualifier": "development"
	//       }
	//     },
	//     "StageName": "development"
	//   },
	//   "Type": "AWS::ApiGateway::Deployment"
	// }
}

func Example_apiDeploymentStaging() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "ApiDeploymentStaging")
	// Output:
	// {
	//   "DependsOn": [
	//     "ApiRootMethod",
	//     "ApiProxyMethod",
	//     "ApiFunctionAliasStaging"
	//   ],
	//   "Properties": {
	//     "RestApiId": {
	//       "Ref": "Api"
	//     },
	//     "StageDescription": {
	//       "Variables": {
	//         "qualifier": "staging"
	//       }
	//     },
	//     "StageName": "staging"
	//   },
	//   "Type": "AWS::ApiGateway::Deployment"
	// }
}

func Example_apiDeploymentProduction() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "ApiDeploymentProduction")
	// Output:
	// {
	//   "DependsOn": [
	//     "ApiRootMethod",
	//     "ApiProxyMethod",
	//     "ApiFunctionAliasProduction"
	//   ],
	//   "Properties": {
	//     "RestApiId": {
	//       "Ref": "Api"
	//     },
	//     "StageDescription": {
	//       "Variables": {
	//         "qualifier": "production"
	//       }
	//     },
	//     "StageName": "production"
	//   },
	//   "Type": "AWS::ApiGateway::Deployment"
	// }
}

func Example_apiFunctionAliasDevelopment() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "ApiFunctionAliasDevelopment")
	// Output:
	// {
	//   "Properties": {
	//     "Description": "Development environment (Managed by Up).",
	//     "FunctionName": {
	//       "Ref": "FunctionName"
	//     },
	//     "FunctionVersion": "$LATEST",
	//     "Name": "development"
	//   },
	//   "Type": "AWS::Lambda::Alias"
	// }
}

func Example_apiFunctionAliasStaging() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "ApiFunctionAliasStaging")
	// Output:
	// {
	//   "Properties": {
	//     "Description": "Staging environment (Managed by Up).",
	//     "FunctionName": {
	//       "Ref": "FunctionName"
	//     },
	//     "FunctionVersion": {
	//       "Ref": "FunctionVersionStaging"
	//     },
	//     "Name": "staging"
	//   },
	//   "Type": "AWS::Lambda::Alias"
	// }
}

func Example_apiFunctionAliasProduction() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
		},
	}

	dump(c, "ApiFunctionAliasProduction")
	// Output:
	// {
	//   "Properties": {
	//     "Description": "Production environment (Managed by Up).",
	//     "FunctionName": {
	//       "Ref": "FunctionName"
	//     },
	//     "FunctionVersion": {
	//       "Ref": "FunctionVersionProduction"
	//     },
	//     "Name": "production"
	//   },
	//   "Type": "AWS::Lambda::Alias"
	// }
}

func Example_apiDomainName() {
	c := &Config{
		Config: up.MustParseConfigString(`{
      "name": "polls",
      "stages": {
        "production": {
          "domain": "gh-polls.com"
        }
      }
    }`),
	}

	dump(c, "ApiDomainProduction")
	// Output:
	// {
	//   "Properties": {
	//     "CertificateArn": "",
	//     "DomainName": "gh-polls.com"
	//   },
	//   "Type": "AWS::ApiGateway::DomainName"
	// }
}

func Example_apiDomainNameMapping() {
	c := &Config{
		Config: up.MustParseConfigString(`{
      "name": "polls",
      "stages": {
        "production": {
          "domain": "gh-polls.com"
        }
      }
    }`),
	}

	dump(c, "ApiDomainProductionPathMapping")
	// Output:
	// {
	//   "DependsOn": [
	//     "ApiDomainProduction",
	//     "ApiDeploymentProduction"
	//   ],
	//   "Properties": {
	//     "BasePath": "",
	//     "DomainName": "gh-polls.com",
	//     "RestApiId": {
	//       "Ref": "Api"
	//     },
	//     "Stage": "production"
	//   },
	//   "Type": "AWS::ApiGateway::BasePathMapping"
	// }
}

func Example_apiDomainDNSRecord() {
	c := &Config{
		Config: up.MustParseConfigString(`{
      "name": "polls",
      "stages": {
        "production": {
          "domain": "gh-polls.com"
        }
      }
    }`),
	}

	dump(c, "DnsZoneGhPollsComRecordGhPollsCom")
	// Output:
	// {
	//   "Properties": {
	//     "AliasTarget": {
	//       "DNSName": {
	//         "Fn::GetAtt": [
	//           "ApiDomainProduction",
	//           "DistributionDomainName"
	//         ]
	//       },
	//       "HostedZoneId": "Z2FDTNDATAQYW2"
	//     },
	//     "Comment": "Managed by Up.",
	//     "HostedZoneId": {
	//       "Ref": "DnsZoneGhPollsCom"
	//     },
	//     "Name": "gh-polls.com",
	//     "Type": "A"
	//   },
	//   "Type": "AWS::Route53::RecordSet"
	// }
}

func Example_apiDomainDNSZone() {
	c := &Config{
		Config: up.MustParseConfigString(`{
      "name": "polls",
      "stages": {
        "production": {
          "domain": "gh-polls.com"
        }
      }
    }`),
	}

	dump(c, "DnsZoneGhPollsCom")
	// Output:
	// {
	//   "Properties": {
	//     "Name": "gh-polls.com"
	//   },
	//   "Type": "AWS::Route53::HostedZone"
	// }
}

func Test_apiDomainDNSZone_existingZoneWithSubdomain(t *testing.T) {
	c := &Config{
		Config: up.MustParseConfigString(`{
      "name": "polls",
      "stages": {
        "production": {
          "domain": "api.gh-polls.com"
        }
      }
    }`),
		Zones: []*route53.HostedZone{
			{
				Name: aws.String("gh-polls.com."),
				Id:   aws.String("SOMETHING"),
			},
		},
	}

	r := New(c)["Resources"].(Map)

	_, ok := r["DnsZoneGhPollsCom"].(Map)
	assert.False(t, ok, "should not create hosted zone")

	record, ok := r["DnsZoneApiGhPollsComRecordApiGhPollsCom"].(Map)
	assert.True(t, ok, "should have A record")

	id := record["Properties"].(Map)["HostedZoneId"]
	assert.Equal(t, "SOMETHING", id)
}

func Example_apiDomainDNSManual() {
	c := &Config{
		Config: up.MustParseConfigString(`{
      "name": "polls",
      "dns": {
        "gh-polls.com": [
          {
            "name": "blog.gh-polls.com",
            "type": "CNAME",
            "value": ["medium.com"]
          }
        ]
      }
    }`),
	}

	dump(c, "DnsZoneGhPollsComRecordBlogGhPollsComCNAME")
	// Output:
	// {
	//   "Properties": {
	//     "Comment": "Managed by Up.",
	//     "HostedZoneId": {
	//       "Ref": "DnsZoneGhPollsCom"
	//     },
	//     "Name": "blog.gh-polls.com",
	//     "ResourceRecords": [
	//       "medium.com"
	//     ],
	//     "TTL": "300",
	//     "Type": "CNAME"
	//   },
	//   "Type": "AWS::Route53::RecordSet"
	// }
}

func Example_iamLambdaDevelopment() {
	c := &Config{
		Config: up.MustParseConfigString(`{
      "name": "polls"
    }`),
	}

	dump(c, "ApiLambdaPermissionDevelopment")
	// Output:
	// {
	//   "DependsOn": "ApiFunctionAliasDevelopment",
	//   "Properties": {
	//     "Action": "lambda:invokeFunction",
	//     "FunctionName": {
	//       "Fn::Join": [
	//         ":",
	//         [
	//           "arn",
	//           "aws",
	//           "lambda",
	//           {
	//             "Ref": "AWS::Region"
	//           },
	//           {
	//             "Ref": "AWS::AccountId"
	//           },
	//           "function",
	//           {
	//             "Fn::Join": [
	//               ":",
	//               [
	//                 {
	//                   "Ref": "FunctionName"
	//                 },
	//                 "development"
	//               ]
	//             ]
	//           }
	//         ]
	//       ]
	//     },
	//     "Principal": "apigateway.amazonaws.com",
	//     "SourceArn": {
	//       "Fn::Join": [
	//         "",
	//         [
	//           "arn:aws:execute-api",
	//           ":",
	//           {
	//             "Ref": "AWS::Region"
	//           },
	//           ":",
	//           {
	//             "Ref": "AWS::AccountId"
	//           },
	//           ":",
	//           {
	//             "Ref": "Api"
	//           },
	//           "/*"
	//         ]
	//       ]
	//     }
	//   },
	//   "Type": "AWS::Lambda::Permission"
	// }
}

func Example_iamLambdaStaging() {
	c := &Config{
		Config: up.MustParseConfigString(`{
      "name": "polls"
    }`),
	}

	dump(c, "ApiLambdaPermissionStaging")
	// Output:
	// {
	//   "DependsOn": "ApiFunctionAliasStaging",
	//   "Properties": {
	//     "Action": "lambda:invokeFunction",
	//     "FunctionName": {
	//       "Fn::Join": [
	//         ":",
	//         [
	//           "arn",
	//           "aws",
	//           "lambda",
	//           {
	//             "Ref": "AWS::Region"
	//           },
	//           {
	//             "Ref": "AWS::AccountId"
	//           },
	//           "function",
	//           {
	//             "Fn::Join": [
	//               ":",
	//               [
	//                 {
	//                   "Ref": "FunctionName"
	//                 },
	//                 "staging"
	//               ]
	//             ]
	//           }
	//         ]
	//       ]
	//     },
	//     "Principal": "apigateway.amazonaws.com",
	//     "SourceArn": {
	//       "Fn::Join": [
	//         "",
	//         [
	//           "arn:aws:execute-api",
	//           ":",
	//           {
	//             "Ref": "AWS::Region"
	//           },
	//           ":",
	//           {
	//             "Ref": "AWS::AccountId"
	//           },
	//           ":",
	//           {
	//             "Ref": "Api"
	//           },
	//           "/*"
	//         ]
	//       ]
	//     }
	//   },
	//   "Type": "AWS::Lambda::Permission"
	// }
}

func Example_iamLambdaProduction() {
	c := &Config{
		Config: up.MustParseConfigString(`{
      "name": "polls"
    }`),
	}

	dump(c, "ApiLambdaPermissionProduction")
	// Output:
	// {
	//   "DependsOn": "ApiFunctionAliasProduction",
	//   "Properties": {
	//     "Action": "lambda:invokeFunction",
	//     "FunctionName": {
	//       "Fn::Join": [
	//         ":",
	//         [
	//           "arn",
	//           "aws",
	//           "lambda",
	//           {
	//             "Ref": "AWS::Region"
	//           },
	//           {
	//             "Ref": "AWS::AccountId"
	//           },
	//           "function",
	//           {
	//             "Fn::Join": [
	//               ":",
	//               [
	//                 {
	//                   "Ref": "FunctionName"
	//                 },
	//                 "production"
	//               ]
	//             ]
	//           }
	//         ]
	//       ]
	//     },
	//     "Principal": "apigateway.amazonaws.com",
	//     "SourceArn": {
	//       "Fn::Join": [
	//         "",
	//         [
	//           "arn:aws:execute-api",
	//           ":",
	//           {
	//             "Ref": "AWS::Region"
	//           },
	//           ":",
	//           {
	//             "Ref": "AWS::AccountId"
	//           },
	//           ":",
	//           {
	//             "Ref": "Api"
	//           },
	//           "/*"
	//         ]
	//       ]
	//     }
	//   },
	//   "Type": "AWS::Lambda::Permission"
	// }
}
