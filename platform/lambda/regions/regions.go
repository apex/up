package regions

import "path/filepath"

// All regions.
var All = []string{
	"us-east-1",
	"us-west-2",
	"eu-west-1",
	"eu-central-1",
	"ap-northeast-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"us-east-2",
	"us-west-1",
	"ap-northeast-2",
	"ap-south-1",
	"sa-east-1",
	"ca-central-1",
}

// Match returns regions matching the pattern(s) provided. Any
// patern which does not match is returned so it may be validated
// (aka it will fail validation).
func Match(regions []string) (v []string) {
	for _, pattern := range regions {
		matched := false

		for _, id := range All {
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
