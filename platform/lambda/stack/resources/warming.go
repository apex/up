package resources

import (
	"fmt"

	"github.com/apex/up/config"
	"github.com/apex/up/internal/util"
)

// warmingFunctionSource is the source code.
var warmingFunctionSource = `
const http = require('https')

exports.handle = function(e, ctx, fn) {
  const start = Date.now()
  let pending = e.count
  console.log('requesting %d', e.count)

  for (let i = 0; i < e.count; i++) {
    console.log('GET %s', e.url)
    http.get(e.url, function (res) {
      const d = Date.now() - start
      console.log('GET %s -> %s (%dms)', e.url, res.statusCode, d)
      --pending || fn()
    })
  }
}`

// warming resources.
func warming(c *Config, m Map) {
	if !c.Lambda.Warm {
		return
	}

	warmingFunctionRole(c, m)
	warmingFunction(c, m)
	warmingFunctionPermission(c, m)
	warmingEvent(c, m)
}

// warmingFunction sets up a scheduled function for warming.
func warmingFunction(c *Config, m Map) {
	m["WarmingFunction"] = Map{
		"Type": "AWS::Lambda::Function",
		"Properties": Map{
			"FunctionName": fmt.Sprintf("%s-warming", c.Name),
			"Description":  util.ManagedByUp("Warming function"),
			"Runtime":      "nodejs6.10",
			"Handler":      "index.handle",
			"Role":         get("WarmingFunctionRole", "Arn"),
			"MemorySize":   512,
			"Timeout":      300,
			"Code": Map{
				"ZipFile": warmingFunctionSource,
			},
		},
	}
}

// warmingFunctionPermission sets up function permissions.
func warmingFunctionPermission(c *Config, m Map) {
	m["WarmingFunctionPermission"] = Map{
		"Type": "AWS::Lambda::Permission",
		"Properties": Map{
			"FunctionName": ref("WarmingFunction"),
			"Action":       "lambda:InvokeFunction",
			"Principal":    "events.amazonaws.com",
			"SourceArn":    get("WarmingEvent", "Arn"),
		},
	}
}

// warmingEvent sets up a warming scheduled event.
func warmingEvent(c *Config, m Map) {
	url := endpoint("production")
	input := join("", `{ "url": "`, url, fmt.Sprintf(`", "count": %d }`, c.Lambda.WarmCount))

	m["WarmingEvent"] = Map{
		"Type": "AWS::Events::Rule",
		"Properties": Map{
			"State":              "ENABLED",
			"Description":        util.ManagedByUp("Warming function scheduled event"),
			"ScheduleExpression": rate(c.Lambda.WarmRate),
			"Targets": []Map{
				{
					"Arn":   get("WarmingFunction", "Arn"),
					"Id":    "WarmingFunction",
					"Input": input,
				},
			},
		},
	}
}

// warmingFunctionRole sets up the warming function role.
func warmingFunctionRole(c *Config, m Map) {
	m["WarmingFunctionRole"] = Map{
		"Type": "AWS::IAM::Role",
		"Properties": Map{
			"RoleName": fmt.Sprintf("%s-warming-function", c.Name),
			"AssumeRolePolicyDocument": Map{
				"Version": "2012-10-17",
				"Statement": []Map{
					{
						"Effect": "Allow",
						"Principal": Map{
							"Service": []string{"lambda.amazonaws.com"},
						},
						"Action": []string{"sts:AssumeRole"},
					},
				},
			},
			"Path": "/",
			"Policies": []Map{
				{
					"PolicyName": "root",
					"PolicyDocument": Map{
						"Version": "2012-10-17",
						"Statement": []Map{
							{
								"Effect":   "Allow",
								"Action":   []string{"logs:*"},
								"Resource": "arn:aws:logs:*:*:*",
							},
						},
					},
				},
			},
		},
	}
}

// rate returns a rate string.
func rate(d config.Duration) string {
	switch m := d.Seconds() / 60; {
	case m == 1:
		return "rate(1 minute)"
	default:
		return fmt.Sprintf("rate(%0.0f minutes)", m)
	}
}

// endpoint returns the api endpoint for stage.
func endpoint(stage string) Map {
	path := fmt.Sprintf("/%s/_ping", stage)
	return join("", "https://", ref("Api"), ".execute-api.", ref("AWS::Region"), ".amazonaws.com", path)
}
