package resources

import (
	"fmt"
	"strconv"

	"github.com/apex/up"
	"github.com/apex/up/config"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/aws/regions"
	"github.com/aws/aws-sdk-go/service/route53"
)

// Map type.
type Map map[string]interface{}

// Versions is a map of stage to lambda function version.
type Versions map[string]string

// Config for the resource template.
type Config struct {
	// Zones already present in route53. This is used to
	// ensure that existing zones previously set up, or
	// automatically configured when purchasing a domain
	// are not duplicated.
	Zones []*route53.HostedZone

	// Versions map used to maintain the correct lambda
	// function aliases when updating a stack.
	Versions Versions

	*up.Config
}

// New template.
func New(c *Config) map[string]interface{} {
	return Map{
		"AWSTemplateFormatVersion": "2010-09-09",
		"Parameters":               parameters(c),
		"Outputs":                  outputs(c),
		"Resources":                resources(c),
	}
}

// ref of id.
func ref(id string) Map {
	return Map{
		"Ref": id,
	}
}

// get value from named ref.
func get(name, value string) Map {
	return Map{
		"Fn::GetAtt": []string{
			name,
			value,
		},
	}
}

// join strings with delim.
func join(delim string, s ...interface{}) Map {
	return Map{
		"Fn::Join": []interface{}{
			delim,
			s,
		},
	}
}

// stageVariable by name.
func stageVariable(name string) string {
	return fmt.Sprintf("${stageVariables.%s}", name)
}

// lambda ARN for function name.
func lambdaArn(name string) Map {
	return join(":", "arn", "aws", "lambda", ref("AWS::Region"), ref("AWS::AccountId"), "function", ref(name))
}

// lambda ARN for function name with qualifier.
func lambdaArnQualifier(name, qualifier string) Map {
	return join(":", "arn", "aws", "lambda", ref("AWS::Region"), ref("AWS::AccountId"), "function", join(":", ref(name), qualifier))
}

// endpointConfiguration returns the endpoint for EDGE or REGIONAL.
func endpointConfiguration(c *Config) Map {
	switch c.Lambda.Endpoint {
	case "regional":
		return Map{
			"Types": []string{"REGIONAL"},
		}
	default:
		return Map{
			"Types": []string{"EDGE"},
		}
	}
}

// getZone returns a zone by domain or nil.
func getZone(c *Config, domain string) *route53.HostedZone {
	for _, z := range c.Zones {
		if *z.Name == domain+"." {
			return z
		}
	}
	return nil
}

// dnsZone returns the ref to a new zone, or id to an existing zone.
func dnsZone(c *Config, m Map, domain string) interface{} {
	// already exists
	if z := getZone(c, domain); z != nil {
		return *z.Id
	}

	id := util.Camelcase("dns_zone_%s", domain)

	// already registered for creation
	if m[id] != nil {
		return ref(id)
	}

	// new zone
	m[id] = Map{
		"Type":                "AWS::Route53::HostedZone",
		"DeletionPolicy":      "Retain",
		"UpdateReplacePolicy": "Retain",
		"Properties": Map{
			"Name": domain,
		},
	}

	return ref(id)
}

// api sets up the app resources.
func api(c *Config, m Map) {
	m["Api"] = Map{
		"Type": "AWS::ApiGateway::RestApi",
		"Properties": Map{
			"Name":        ref("Name"),
			"Description": util.ManagedByUp(c.Description),
			"BinaryMediaTypes": []string{
				"*/*",
			},
			"EndpointConfiguration": endpointConfiguration(c),
		},
	}

	integration := Map{
		"Type":                  "AWS_PROXY",
		"IntegrationHttpMethod": "POST",
		"Uri": join("",
			"arn:aws:apigateway:",
			ref("AWS::Region"),
			":lambda:path/2015-03-31/functions/",
			lambdaArnQualifier("FunctionName", stageVariable("qualifier")),
			"/invocations"),
	}

	m["ApiRootMethod"] = Map{
		"Type": "AWS::ApiGateway::Method",
		"Properties": Map{
			"RestApiId":         ref("Api"),
			"ResourceId":        get("Api", "RootResourceId"),
			"HttpMethod":        "ANY",
			"AuthorizationType": "NONE",
			"Integration":       integration,
		},
	}

	m["ApiProxyResource"] = Map{
		"Type": "AWS::ApiGateway::Resource",
		"Properties": Map{
			"RestApiId": ref("Api"),
			"ParentId":  get("Api", "RootResourceId"),
			"PathPart":  "{proxy+}",
		},
	}

	m["ApiProxyMethod"] = Map{
		"Type": "AWS::ApiGateway::Method",
		"Properties": Map{
			"RestApiId":         ref("Api"),
			"ResourceId":        ref("ApiProxyResource"),
			"HttpMethod":        "ANY",
			"AuthorizationType": "NONE",
			"Integration":       integration,
		},
	}

	stages(c, m)
}

// stages sets up the stage specific resources.
func stages(c *Config, m Map) {
	for _, s := range c.Stages.List() {
		if s.IsRemote() {
			stage(c, s, m)
		}
	}
}

// stage sets up the stage specific resources.
func stage(c *Config, s *config.Stage, m Map) {
	aliasID := stageAlias(c, s, m)
	deploymentID := stageDeployment(c, s, m, aliasID)
	stagePermissions(c, s, m, aliasID)
	stageDomain(c, s, m, deploymentID)
}

// stageAlias sets up the lambda alias and deployment and returns the alias id.
func stageAlias(c *Config, s *config.Stage, m Map) string {
	id := util.Camelcase("api_function_alias_%s", s.Name)
	version, ok := c.Versions[s.Name]

	if !ok {
		panic(fmt.Sprintf("stage %q is missing a function version mapping", s.Name))
	}

	m[id] = Map{
		"Type": "AWS::Lambda::Alias",
		"Properties": Map{
			"Name":            s.Name,
			"Description":     util.ManagedByUp(""),
			"FunctionName":    ref("FunctionName"),
			"FunctionVersion": version,
		},
	}

	return id
}

// stagePermissions sets up the lambda:invokeFunction permissions for API Gateway.
func stagePermissions(c *Config, s *config.Stage, m Map, aliasID string) {
	id := util.Camelcase("api_lambda_permission_%s", s.Name)

	m[id] = Map{
		"Type":      "AWS::Lambda::Permission",
		"DependsOn": aliasID,
		"Properties": Map{
			"Action":       "lambda:invokeFunction",
			"FunctionName": lambdaArnQualifier("FunctionName", s.Name),
			"Principal":    "apigateway.amazonaws.com",
			"SourceArn": join("",
				"arn:aws:execute-api",
				":",
				ref("AWS::Region"),
				":",
				ref("AWS::AccountId"),
				":",
				ref("Api"),
				"/*"),
		},
	}
}

// stageDeployment sets up the API Gateway deployment.
func stageDeployment(c *Config, s *config.Stage, m Map, aliasID string) string {
	id := util.Camelcase("api_deployment_%s", s.Name)

	m[id] = Map{
		"Type":      "AWS::ApiGateway::Deployment",
		"DependsOn": []string{"ApiRootMethod", "ApiProxyMethod", aliasID},
		"Properties": Map{
			"RestApiId": ref("Api"),
			"StageName": s.Name,
			"StageDescription": Map{
				"Variables": Map{
					"qualifier": s.Name,
				},
			},
		},
	}

	return id
}

// stageDomain sets up a custom domain, dns record and path mapping.
func stageDomain(c *Config, s *config.Stage, m Map, deploymentID string) {
	if s.Domain == "" {
		return
	}

	id := util.Camelcase("api_domain_%s", s.Name)

	props := Map{
		"DomainName":            s.Domain,
		"EndpointConfiguration": endpointConfiguration(c),
	}

	if c.Lambda.Endpoint == "regional" {
		props["RegionalCertificateArn"] = s.Cert
	} else {
		props["CertificateArn"] = s.Cert

	}

	m[id] = Map{
		"Type":       "AWS::ApiGateway::DomainName",
		"Properties": props,
	}

	stagePathMapping(c, s, m, deploymentID, id)

	if s.Zone != false {
		stageDNSRecord(c, s, m, id)
	}
}

// stagePathMapping sets up the stage deployment mapping.
func stagePathMapping(c *Config, s *config.Stage, m Map, deploymentID, domainID string) {
	id := util.Camelcase("api_domain_%s_path_mapping", s.Name)

	m[id] = Map{
		"Type":      "AWS::ApiGateway::BasePathMapping",
		"DependsOn": []string{deploymentID, domainID},
		"Properties": Map{
			"DomainName": s.Domain,
			"BasePath":   util.BasePath(s.Path),
			"RestApiId":  ref("Api"),
			"Stage":      s.Name,
		},
	}
}

// stageDNSRecord sets up an ALIAS record and zone if necessary for a custom domain.
func stageDNSRecord(c *Config, s *config.Stage, m Map, domainID string) {
	if c.Lambda.Endpoint == "regional" {
		stageDNSRecordRegional(c, s, m, domainID)
	} else {
		stageDNSRecordEdge(c, s, m, domainID)
	}
}

// stageDNSRecordRegional sets up an A record with latency routing policy.
func stageDNSRecordRegional(c *Config, s *config.Stage, m Map, domainID string) {
	id := util.Camelcase("dns_zone_%s_record_%s", util.Domain(s.Domain), s.Domain)
	zoneName := util.Domain(s.Domain)

	// explicit .zone was specified
	if s, ok := s.Zone.(string); ok {
		zoneName = s
	}

	zone := dnsZone(c, m, zoneName)
	region := c.Regions[0]

	// api gateway
	hostedZoneID := regions.GetHostedZoneID(region)

	m[id] = Map{
		"Type": "AWS::Route53::RecordSet",
		"Properties": Map{
			"Name":          s.Domain,
			"Type":          "A",
			"Comment":       util.ManagedByUp(""),
			"Region":        region,
			"SetIdentifier": region,
			"HostedZoneId":  zone,
			"AliasTarget": Map{
				"DNSName":      dnsName(c, domainID),
				"HostedZoneId": hostedZoneID,
			},
		},
	}
}

// stageDNSRecordEdge sets up an A record.
func stageDNSRecordEdge(c *Config, s *config.Stage, m Map, domainID string) {
	id := util.Camelcase("dns_zone_%s_record_%s", util.Domain(s.Domain), s.Domain)
	zoneName := util.Domain(s.Domain)

	// explicit .zone was specified
	if s, ok := s.Zone.(string); ok {
		zoneName = s
	}

	zone := dnsZone(c, m, zoneName)

	m[id] = Map{
		"Type": "AWS::Route53::RecordSet",
		"Properties": Map{
			"Name":         s.Domain,
			"Type":         "A",
			"Comment":      util.ManagedByUp(""),
			"HostedZoneId": zone,
			"AliasTarget": Map{
				"DNSName":      dnsName(c, domainID),
				"HostedZoneId": "Z2FDTNDATAQYW2", // Cloudfront
			},
		},
	}
}

// dnsName returns the endpoint based on the Lambda.Endpoint configuration.
func dnsName(c *Config, domainID string) Map {
	if c.Lambda.Endpoint == "regional" {
		return get(domainID, "RegionalDomainName")
	}
	return get(domainID, "DistributionDomainName")
}

// dns setups the the user-defined DNS zones and records.
func dns(c *Config, m Map) {
	for _, z := range c.DNS.Zones {
		zone := dnsZone(c, m, z.Name)

		for _, r := range z.Records {
			id := util.Camelcase("dns_zone_%s_record_%s_%s", z.Name, r.Name, r.Type)

			m[id] = Map{
				"Type": "AWS::Route53::RecordSet",
				"Properties": Map{
					"Name":            r.Name,
					"Type":            r.Type,
					"TTL":             strconv.Itoa(r.TTL),
					"ResourceRecords": r.Value,
					"HostedZoneId":    zone,
					"Comment":         util.ManagedByUp(""),
				},
			}
		}
	}
}

// resources of the stack.
func resources(c *Config) Map {
	m := Map{}
	api(c, m)
	dns(c, m)
	alerting(c, m)
	warming(c, m)
	return m
}

// parameters of the stack.
func parameters(c *Config) Map {
	return Map{
		"Name": Map{
			"Description": "Name of application",
			"Type":        "String",
		},
		"FunctionName": Map{
			"Description": "Name of application function",
			"Type":        "String",
		},
	}
}

// outputs of the stack.
func outputs(c *Config) Map {
	return Map{
		"ApiName": Map{
			"Description": "API name",
			"Value":       ref("Name"),
		},
		"ApiFunctionName": Map{
			"Description": "API Lambda function name",
			"Value":       ref("FunctionName"),
		},
		"ApiFunctionArn": Map{
			"Description": "API Lambda function ARN",
			"Value":       lambdaArn("FunctionName"),
		},
	}
}
