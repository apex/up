package resources

import (
	"time"

	"github.com/apex/up"
	"github.com/apex/up/config"
)

func Example_warmingFunction() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Lambda: config.Lambda{
				Warm: true,
			},
		},
	}

	dump(c, "WarmingFunction")
	// Output:
	// {
	//   "Properties": {
	//     "Code": {
	//       "ZipFile": "\nconst http = require('https')\n\nexports.handle = function(e, ctx, fn) {\n  const start = Date.now()\n  let pending = e.count\n  console.log('requesting %d', e.count)\n\n  for (let i = 0; i \u003c e.count; i++) {\n    console.log('GET %s', e.url)\n    http.get(e.url, function (res) {\n      const d = Date.now() - start\n      console.log('GET %s -\u003e %s (%dms)', e.url, res.statusCode, d)\n      --pending || fn()\n    })\n  }\n}"
	//     },
	//     "Description": "Warming function (Managed by Up).",
	//     "FunctionName": "polls-warming",
	//     "Handler": "index.handle",
	//     "MemorySize": 512,
	//     "Role": {
	//       "Fn::GetAtt": [
	//         "WarmingFunctionRole",
	//         "Arn"
	//       ]
	//     },
	//     "Runtime": "nodejs6.10",
	//     "Timeout": 300
	//   },
	//   "Type": "AWS::Lambda::Function"
	// }
}

func Example_warmingFunctionPermission() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Lambda: config.Lambda{
				Warm: true,
			},
		},
	}

	dump(c, "WarmingFunctionPermission")
	// Output:
	// {
	//   "Properties": {
	//     "Action": "lambda:InvokeFunction",
	//     "FunctionName": {
	//       "Ref": "WarmingFunction"
	//     },
	//     "Principal": "events.amazonaws.com",
	//     "SourceArn": {
	//       "Fn::GetAtt": [
	//         "WarmingEvent",
	//         "Arn"
	//       ]
	//     }
	//   },
	//   "Type": "AWS::Lambda::Permission"
	// }
}

func Example_warmingEvent() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Lambda: config.Lambda{
				Warm:      true,
				WarmRate:  config.Duration(5 * time.Minute),
				WarmCount: 30,
			},
		},
	}

	dump(c, "WarmingEvent")
	// Output:
	// {
	//   "Properties": {
	//     "Description": "Warming function scheduled event (Managed by Up).",
	//     "ScheduleExpression": "rate(5 minutes)",
	//     "State": "ENABLED",
	//     "Targets": [
	//       {
	//         "Arn": {
	//           "Fn::GetAtt": [
	//             "WarmingFunction",
	//             "Arn"
	//           ]
	//         },
	//         "Id": "WarmingFunction",
	//         "Input": {
	//           "Fn::Join": [
	//             "",
	//             [
	//               "{ \"url\": \"",
	//               {
	//                 "Fn::Join": [
	//                   "",
	//                   [
	//                     "https://",
	//                     {
	//                       "Ref": "Api"
	//                     },
	//                     ".execute-api.",
	//                     {
	//                       "Ref": "AWS::Region"
	//                     },
	//                     ".amazonaws.com",
	//                     "/production/_ping"
	//                   ]
	//                 ]
	//               },
	//               "\", \"count\": 30 }"
	//             ]
	//           ]
	//         }
	//       }
	//     ]
	//   },
	//   "Type": "AWS::Events::Rule"
	// }
}

func Example_warmingFunctionRole() {
	c := &Config{
		Config: &up.Config{
			Name: "polls",
			Lambda: config.Lambda{
				Warm:      true,
				WarmRate:  config.Duration(5 * time.Minute),
				WarmCount: 30,
			},
		},
	}

	dump(c, "WarmingFunctionRole")
	// Output:
	// {
	//   "Properties": {
	//     "AssumeRolePolicyDocument": {
	//       "Statement": [
	//         {
	//           "Action": [
	//             "sts:AssumeRole"
	//           ],
	//           "Effect": "Allow",
	//           "Principal": {
	//             "Service": [
	//               "lambda.amazonaws.com"
	//             ]
	//           }
	//         }
	//       ],
	//       "Version": "2012-10-17"
	//     },
	//     "Path": "/",
	//     "Policies": [
	//       {
	//         "PolicyDocument": {
	//           "Statement": [
	//             {
	//               "Action": [
	//                 "logs:*"
	//               ],
	//               "Effect": "Allow",
	//               "Resource": "arn:aws:logs:*:*:*"
	//             }
	//           ],
	//           "Version": "2012-10-17"
	//         },
	//         "PolicyName": "root"
	//       }
	//     ],
	//     "RoleName": "polls-warming-function"
	//   },
	//   "Type": "AWS::IAM::Role"
	// }
}
