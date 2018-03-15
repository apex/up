// Package regions provides AWS region utilities.
package regions

import "path/filepath"

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
