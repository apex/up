package config

import (
	"time"

	"github.com/pkg/errors"

	"github.com/apex/up/internal/validate"
)

// defaultRuntime is the default runtime.
var defaultRuntime = "nodejs10.x"

// defaultPolicy is the default function role policy.
var defaultPolicy = IAMPolicyStatement{
	"Effect":   "Allow",
	"Resource": "*",
	"Action": []string{
		"logs:CreateLogGroup",
		"logs:CreateLogStream",
		"logs:PutLogEvents",
		"ssm:GetParametersByPath",
		"ec2:CreateNetworkInterface",
		"ec2:DescribeNetworkInterfaces",
		"ec2:DeleteNetworkInterface",
	},
}

// IAMPolicyStatement configuration.
type IAMPolicyStatement map[string]interface{}

// VPC configuration.
type VPC struct {
	Subnets        []string `json:"subnets"`
	SecurityGroups []string `json:"security_groups"`
}

// Lambda configuration.
type Lambda struct {
	// Memory of the function.
	Memory int `json:"memory"`

	// Timeout of the function.
	Timeout int `json:"timeout"`

	// Role of the function.
	Role string `json:"role"`

	// Runtime of the function.
	Runtime string `json:"runtime"`

	// Policy of the function role.
	Policy []IAMPolicyStatement `json:"policy"`

	// VPC configuration.
	VPC *VPC `json:"vpc"`

	// Accelerate enables S3 acceleration.
	Accelerate bool `json:"accelerate"`

	// Warm enables active warming.
	Warm *bool `json:"warm"`

	// WarmCount is the number of containers to keep active.
	WarmCount int `json:"warm_count"`

	// WarmRate is the schedule for performing worming.
	WarmRate Duration `json:"warm_rate"`

	// Layers to include.
	Layers []string `json:"layers"`

	// Endpoint type, "edge" or "regional", defaults to "edge".
	Endpoint string `json:"endpoint"`
}

// Default implementation.
func (l *Lambda) Default() error {
	if l.Timeout == 0 {
		l.Timeout = 60
	}

	if l.Memory == 0 {
		l.Memory = 512
	}

	if l.Runtime == "" {
		l.Runtime = defaultRuntime
	}

	l.Policy = append(l.Policy, defaultPolicy)

	if l.WarmRate == 0 {
		l.WarmRate = Duration(15 * time.Minute)
	}

	if l.WarmCount == 0 {
		l.WarmCount = 15
	}

	if l.Endpoint == "" {
		l.Endpoint = "edge"
	}

	return nil
}

// Validate implementation.
func (l *Lambda) Validate() error {
	if err := validate.List(l.Endpoint, []string{"edge", "regional"}); err != nil {
		return errors.Wrap(err, ".endpoint")
	}

	if err := validate.List(l.Endpoint, []string{"edge", "regional"}); err != nil {
		return errors.Wrap(err, ".endpoint")
	}

	return nil
}

// Override config.
func (l *Lambda) Override(c *Config) {
	if l.Memory != 0 {
		c.Lambda.Memory = l.Memory
	}

	if l.Timeout != 0 {
		c.Lambda.Timeout = l.Timeout
	}

	if l.Role != "" {
		c.Lambda.Role = l.Role
	}

	if l.VPC != nil {
		c.Lambda.VPC = l.VPC
	}

	if l.Runtime != "" {
		c.Lambda.Runtime = l.Runtime
	}

	if l.Warm != nil {
		c.Lambda.Warm = l.Warm
	}

	if l.WarmCount > 0 {
		c.Lambda.WarmCount = l.WarmCount
	}

	if l.WarmRate != 0 {
		c.Lambda.WarmRate = l.WarmRate
	}

	if l.Endpoint != "" {
		c.Lambda.Endpoint = l.Endpoint
	}
}
