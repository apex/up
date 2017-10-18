package reporter

import "strings"

// TODO: move most of reporting here

// types map.
var types = map[string]string{
	"AWS::Lambda::Alias":               "Lambda alias",
	"AWS::Lambda::Permission":          "Lambda permission",
	"AWS::ApiGateway::RestApi":         "API",
	"AWS::ApiGateway::Method":          "API method",
	"AWS::ApiGateway::Deployment":      "API deployment",
	"AWS::ApiGateway::Resource":        "API resource",
	"AWS::ApiGateway::DomainName":      "API domain",
	"AWS::ApiGateway::BasePathMapping": "API mapping",
	"AWS::Route53::HostedZone":         "DNS zone",
	"AWS::Route53::RecordSet":          "DNS record",
}

// ResourceType returns a human-friendly resource type name.
func ResourceType(s string) string {
	if types[s] != "" {
		return strings.ToLower(types[s])
	}

	return s
}
