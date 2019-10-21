package resources

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/apex/up/config"
	"github.com/tj/assert"
)

func parse(s string) *Config {
	return &Config{
		Config: config.MustParseConfigString(s),
		Versions: Versions{
			"beta":       "1",
			"staging":    "1",
			"production": "1",
		},
	}
}

func assertJSON(t testing.TB, expected string, actual interface{}) {
	t.Helper()
	b, err := json.MarshalIndent(actual, "", "  ")
	assert.NoError(t, err, "marshal")
	s := string(b)
	if expected != s {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n\n", expected, s)
	}
}

var functionJSON = `{
  "Properties": {
    "Code": {
      "ZipFile": "\nconst http = require('https')\n\nexports.handle = function(e, ctx, fn) {\n  const start = Date.now()\n  let pending = e.count\n  console.log('requesting %d', e.count)\n\n  for (let i = 0; i \u003c e.count; i++) {\n    console.log('GET %s', e.url)\n    http.get(e.url, function (res) {\n      const d = Date.now() - start\n      console.log('GET %s -\u003e %s (%dms)', e.url, res.statusCode, d)\n      --pending || fn()\n    })\n  }\n}"
    },
    "Description": "Warming function (Managed by Up).",
    "FunctionName": "polls-warming",
    "Handler": "index.handle",
    "MemorySize": 512,
    "Role": {
      "Fn::GetAtt": [
        "WarmingFunctionRole",
        "Arn"
      ]
    },
    "Runtime": "nodejs10.x",
    "Timeout": 300
  },
  "Type": "AWS::Lambda::Function"
}`

var roleJSON = `{
  "Properties": {
    "AssumeRolePolicyDocument": {
      "Statement": [
        {
          "Action": [
            "sts:AssumeRole"
          ],
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "lambda.amazonaws.com"
            ]
          }
        }
      ],
      "Version": "2012-10-17"
    },
    "Path": "/",
    "Policies": [
      {
        "PolicyDocument": {
          "Statement": [
            {
              "Action": [
                "logs:*"
              ],
              "Effect": "Allow",
              "Resource": "arn:aws:logs:*:*:*"
            }
          ],
          "Version": "2012-10-17"
        },
        "PolicyName": "root"
      }
    ],
    "RoleName": "polls-warming-function"
  },
  "Type": "AWS::IAM::Role"
}`

var eventJSON = `{
  "Properties": {
    "Description": "Warming function scheduled event (Managed by Up).",
    "ScheduleExpression": "rate(15 minutes)",
    "State": "ENABLED",
    "Targets": [
      {
        "Arn": {
          "Fn::GetAtt": [
            "WarmingFunction",
            "Arn"
          ]
        },
        "Id": "WarmingFunction",
        "Input": {
          "Fn::Join": [
            "",
            [
              "{ \"url\": \"",
              {
                "Fn::Join": [
                  "",
                  [
                    "https://",
                    {
                      "Ref": "Api"
                    },
                    ".execute-api.",
                    {
                      "Ref": "AWS::Region"
                    },
                    ".amazonaws.com",
                    "/production/_ping"
                  ]
                ]
              },
              "\", \"count\": 15 }"
            ]
          ]
        }
      }
    ]
  },
  "Type": "AWS::Events::Rule"
}`

var permissionJSON = `{
  "Properties": {
    "Action": "lambda:InvokeFunction",
    "FunctionName": {
      "Ref": "WarmingFunction"
    },
    "Principal": "events.amazonaws.com",
    "SourceArn": {
      "Fn::GetAtt": [
        "WarmingEventProduction",
        "Arn"
      ]
    }
  },
  "Type": "AWS::Lambda::Permission"
}`

func TestWarming_none(t *testing.T) {
	c := parse(`{
		"name": "polls"
	}`)

	assert.Nil(t, getResource(c, "WarmingFunction"))
	assert.Nil(t, getResource(c, "WarmingFunctionRole"))
}

func TestWarming_globalDefaults(t *testing.T) {
	c := parse(`{
		"name": "polls",
		"lambda": {
			"warm": true
		}
	}`)

	assertJSON(t, functionJSON, getResource(c, "WarmingFunction"))
	assertJSON(t, roleJSON, getResource(c, "WarmingFunctionRole"))

	assertJSON(t, strings.Replace(eventJSON, "production", "staging", 1), getResource(c, "WarmingEventStaging"))
	assertJSON(t, eventJSON, getResource(c, "WarmingEventProduction"))

	assertJSON(t, strings.Replace(permissionJSON, "Production", "Staging", 1), getResource(c, "WarmingFunctionPermissionStaging"))
	assertJSON(t, permissionJSON, getResource(c, "WarmingFunctionPermissionProduction"))
}

func TestWarming_globalValues(t *testing.T) {
	c := parse(`{
		"name": "polls",
		"lambda": {
			"warm": true,
			"warm_count": 60,
			"warm_rate": "5m"
		},
		"stages": {
			"beta": {}
		}
	}`)

	assertJSON(t, functionJSON, getResource(c, "WarmingFunction"))
	assertJSON(t, roleJSON, getResource(c, "WarmingFunctionRole"))

	e := strings.Replace(eventJSON, ": 15", ": 60", 1)
	e = strings.Replace(e, "rate(15 minutes)", "rate(5 minutes)", 1)

	assertJSON(t, strings.Replace(e, "production", "staging", 1), getResource(c, "WarmingEventStaging"))
	assertJSON(t, strings.Replace(e, "production", "beta", 1), getResource(c, "WarmingEventBeta"))
	assertJSON(t, e, getResource(c, "WarmingEventProduction"))

	assertJSON(t, strings.Replace(permissionJSON, "Production", "Staging", 1), getResource(c, "WarmingFunctionPermissionStaging"))
	assertJSON(t, strings.Replace(permissionJSON, "Production", "Beta", 1), getResource(c, "WarmingFunctionPermissionBeta"))
	assertJSON(t, permissionJSON, getResource(c, "WarmingFunctionPermissionProduction"))
}

func TestWarming_globalValuesWithOverrides(t *testing.T) {
	c := parse(`{
		"name": "polls",
		"lambda": {
			"warm": true,
			"warm_count": 60,
			"warm_rate": "5m"
		},
		"stages": {
			"beta": {
				"lambda": {
					"warm": false
				}
			},
			"production": {
				"lambda": {
					"warm_count": 100,
					"warm_rate": "10m"
				}
			}
		}
	}`)

	// TODO: test for warm count override
	// TODO: test for warm rate override
	// TODO: it is a pointer problem...

	assertJSON(t, functionJSON, getResource(c, "WarmingFunction"))
	assertJSON(t, roleJSON, getResource(c, "WarmingFunctionRole"))
	assert.Nil(t, getResource(c, "WarmingEventBeta"), "beta disabled")

	prod := strings.Replace(eventJSON, ": 15", ": 100", 1)
	prod = strings.Replace(prod, "rate(15 minutes)", "rate(10 minutes)", 1)
	assertJSON(t, prod, getResource(c, "WarmingEventProduction"))

	stage := strings.Replace(eventJSON, ": 15", ": 60", 1)
	stage = strings.Replace(stage, "rate(15 minutes)", "rate(5 minutes)", 1)
	stage = strings.Replace(stage, "production", "staging", 1)
	assertJSON(t, stage, getResource(c, "WarmingEventStaging"))

	assertJSON(t, strings.Replace(permissionJSON, "Production", "Staging", 1), getResource(c, "WarmingFunctionPermissionStaging"))
	assert.Nil(t, getResource(c, "WarmingFunctionPermissionBeta"), "beta disabled")
	assertJSON(t, permissionJSON, getResource(c, "WarmingFunctionPermissionProduction"))
}
