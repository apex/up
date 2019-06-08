package config

import (
	"errors"
)

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
}

// Default implementation.
func (l *Lambda) Default() error {
	if l.Memory == 0 {
		l.Memory = 512
	}

	if l.Runtime == "" {
		l.Runtime = "nodejs10.x"
	}

	l.Policy = append(l.Policy, defaultPolicy)

	return nil
}

// Validate implementation.
func (l *Lambda) Validate() error {
	if l.Timeout != 0 {
		return errors.New(".lambda.timeout is deprecated, use .proxy.timeout")
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
}
