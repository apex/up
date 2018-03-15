package resources

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/apex/up"
	"github.com/apex/up/config"
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
	v, _ := r[name].(Map)
	return v
}

// dump a resource to stdout.
func dump(c *Config, name string) {
	tmpl := New(c)
	r := tmpl["Resources"].(Map)

	v, ok := r[name].(Map)
	if !ok {
		k := strings.Join(keys(r), "\n  - ")
		panic(fmt.Sprintf("resource %q does not exist in:\n\n  - %s", name, k))
	}

	{
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(v)
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

func Example_stageAlias() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Stages: config.Stages{
				"production": &config.Stage{
					Name: "production",
				},
			},
		},
		Versions: Versions{
			"production": "15",
		},
	}

	dump(c, "ApiFunctionAliasProduction")
	// Output:
	// {
	//   "Properties": {
	//     "Description": "Managed by Up.",
	//     "FunctionName": {
	//       "Ref": "FunctionName"
	//     },
	//     "FunctionVersion": "15",
	//     "Name": "production"
	//   },
	//   "Type": "AWS::Lambda::Alias"
	// }
}

func Example_stagePermission() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Stages: config.Stages{
				"production": &config.Stage{
					Name: "production",
				},
			},
		},
		Versions: Versions{
			"production": "15",
		},
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

func Example_stageDeployment() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Stages: config.Stages{
				"production": &config.Stage{
					Name: "production",
				},
			},
		},
		Versions: Versions{
			"production": "15",
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

func Example_stageDomain() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Stages: config.Stages{
				"production": &config.Stage{
					Name:   "production",
					Domain: "up-example.com",
					Cert:   "arn::something",
				},
			},
		},
		Versions: Versions{
			"production": "15",
		},
	}

	dump(c, "ApiDomainProduction")
	// Output:
	// 	{
	//   "Properties": {
	//     "CertificateArn": "arn::something",
	//     "DomainName": "up-example.com"
	//   },
	//   "Type": "AWS::ApiGateway::DomainName"
	// }
}

func Example_stagePathMapping() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Stages: config.Stages{
				"production": &config.Stage{
					Name:   "production",
					Domain: "up-example.com",
				},
			},
		},
		Versions: Versions{
			"production": "15",
		},
	}

	dump(c, "ApiDomainProductionPathMapping")
	// Output:
	// {
	//   "DependsOn": [
	//     "ApiDeploymentProduction",
	//     "ApiDomainProduction"
	//   ],
	//   "Properties": {
	//     "BasePath": "",
	//     "DomainName": "up-example.com",
	//     "RestApiId": {
	//       "Ref": "Api"
	//     },
	//     "Stage": "production"
	//   },
	//   "Type": "AWS::ApiGateway::BasePathMapping"
	// }
}

func Example_stageDNSZone() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Stages: config.Stages{
				"production": &config.Stage{
					Name:   "production",
					Domain: "up-example.com",
				},
			},
		},
		Versions: Versions{
			"production": "15",
		},
	}

	dump(c, "DnsZoneUpExampleCom")
	// Output:
	// {
	//   "DeletionPolicy": "Retain",
	//   "Properties": {
	//     "Name": "up-example.com"
	//   },
	//   "Type": "AWS::Route53::HostedZone"
	// }
}

func Example_stageDNSZoneRecord() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Stages: config.Stages{
				"production": &config.Stage{
					Name:   "production",
					Domain: "up-example.com",
				},
			},
		},
		Versions: Versions{
			"production": "15",
		},
	}

	dump(c, "DnsZoneUpExampleComRecordUpExampleCom")
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
	//       "Ref": "DnsZoneUpExampleCom"
	//     },
	//     "Name": "up-example.com",
	//     "Type": "A"
	//   },
	//   "Type": "AWS::Route53::RecordSet"
	// }
}

func Example_dnsZone() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			DNS: config.DNS{
				Zones: []*config.Zone{
					{
						Name: "up-example.com",
					},
				},
			},
		},
	}

	dump(c, "DnsZoneUpExampleCom")
	// Output:
	// {
	//   "DeletionPolicy": "Retain",
	//   "Properties": {
	//     "Name": "up-example.com"
	//   },
	//   "Type": "AWS::Route53::HostedZone"
	// }
}

func Example_dnsZoneRecord() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			DNS: config.DNS{
				Zones: []*config.Zone{
					{
						Name: "up-example.com",
						Records: []*config.Record{
							{
								Name:  "blog.up-example.com",
								Type:  "CNAME",
								TTL:   600,
								Value: []string{"example.medium.com"},
							},
						},
					},
				},
			},
		},
	}

	dump(c, "DnsZoneUpExampleComRecordBlogUpExampleComCNAME")
	// Output:
	// {
	//   "Properties": {
	//     "Comment": "Managed by Up.",
	//     "HostedZoneId": {
	//       "Ref": "DnsZoneUpExampleCom"
	//     },
	//     "Name": "blog.up-example.com",
	//     "ResourceRecords": [
	//       "example.medium.com"
	//     ],
	//     "TTL": "600",
	//     "Type": "CNAME"
	//   },
	//   "Type": "AWS::Route53::RecordSet"
	// }
}
