package config

import (
	"bytes"
	"strconv"
	"time"
)

// Duration may be specified as numerical seconds or
// as a duration string such as "1.5m".
type Duration time.Duration

// UnmarshalJSON implementation.
func (d *Duration) UnmarshalJSON(b []byte) error {
	if i, err := strconv.ParseInt(string(b), 10, 64); err == nil {
		*d = Duration(time.Second * time.Duration(i))
		return nil
	}

	v, err := time.ParseDuration(string(bytes.Trim(b, `"`)))
	if err != nil {
		return err
	}

	*d = Duration(v)
	return nil
}
