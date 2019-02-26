// Package regions provides AWS region utilities.
package regions

import (
	"fmt"
	"path/filepath"
)

// hostedZoneIDs is a set of hosted zone ids for API Gateway.
var hostedZoneIDs = map[string]string{
	"us-east-2":      "ZOJJZC49E0EPZ",
	"us-east-1":      "Z1UJRXOUMOOFQ8",
	"us-west-1":      "Z2MUQ32089INYE",
	"us-west-2":      "Z2OJLYMUO9EFXC",
	"ap-south-1":     "Z3VO1THU9YC4UR",
	"ap-northeast-3": "Z2YQB5RD63NC85",
	"ap-northeast-2": "Z20JF4UZKIW1U8",
	"ap-southeast-1": "ZL327KTPIQFUL",
	"ap-southeast-2": "Z2RPCDW04V8134",
	"ap-northeast-1": "Z1YSHQZHG15GKL",
	"ca-central-1":   "Z19DQILCV0OWEC",
	"eu-central-1":   "Z1U9ULNL0V5AJ3",
	"eu-west-1":      "ZLY8HYME6SFDD",
	"eu-west-2":      "ZJ5UAJN8Y3Z2Q",
	"eu-west-3":      "Z3KY65QIEKYHQQ",
	"eu-north-1":     "Z2YB950C88HT6D",
	"sa-east-1":      "ZCMLWB8V5SYIT",
}

// IDs of regions.
var IDs = []string{
	"us-east-2",
	"us-east-1",
	"us-west-1",
	"us-west-2",
	"ap-south-1",
	"ap-northeast-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-northeast-1",
	"ca-central-1",
	"eu-central-1",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"sa-east-1",
}

// Names of regions.
var Names = []string{
	"US East (Ohio)",
	"US East (N. Virginia)",
	"US West (N. California)",
	"US West (Oregon)",
	"Asia Pacific (Mumbai)",
	"Asia Pacific (Seoul)",
	"Asia Pacific (Singapore)",
	"Asia Pacific (Sydney)",
	"Asia Pacific (Tokyo)",
	"Canada (Central)",
	"EU (Frankfurt)",
	"EU (Ireland)",
	"EU (London)",
	"EU (Paris)",
	"South America (São Paulo)",
}

// Match returns regions matching the pattern(s) provided. Patterns
// which are not "expanded" are returned as-is.
func Match(regions []string) (v []string) {
	for _, pattern := range regions {
		matched := false

		for _, id := range IDs {
			if ok, _ := filepath.Match(pattern, id); ok {
				v = append(v, id)
				matched = true
			}
		}

		if !matched {
			v = append(v, pattern)
		}
	}

	return
}

// GetIdByName returns a region id by name.
func GetIdByName(name string) string {
	for i, n := range Names {
		if n == name {
			return IDs[i]
		}
	}
	return ""
}

// GetHostedZoneID returns a hosted zone id by region.
func GetHostedZoneID(region string) string {
	id, ok := hostedZoneIDs[region]
	if !ok {
		panic(fmt.Sprintf("region %q is not yet supported", region))
	}
	return id
}
