package config

import (
	"encoding/json"

	"github.com/apex/up/internal/validate"
	"github.com/pkg/errors"
)

// recordTypes is a list of valid record types.
var recordTypes = []string{
	"ALIAS",
	"A",
	"AAAA",
	"CNAME",
	"MX",
	"NAPTR",
	"NS",
	"PTR",
	"SOA",
	"SPF",
	"SRV",
	"TXT",
}

// DNS config.
type DNS struct {
	Zones []*Zone `json:"zones"`
}

// UnmarshalJSON implementation.
func (d *DNS) UnmarshalJSON(b []byte) error {
	var zones map[string][]*Record

	if err := json.Unmarshal(b, &zones); err != nil {
		return err
	}

	for name, records := range zones {
		zone := &Zone{Name: name}
		d.Zones = append(d.Zones, zone)
		for _, record := range records {
			zone.Records = append(zone.Records, record)
		}
	}

	return nil
}

// Default implementation.
func (d *DNS) Default() error {
	for _, z := range d.Zones {
		if err := z.Default(); err != nil {
			return errors.Wrapf(err, "zone %s", z.Name)
		}
	}

	return nil
}

// Validate implementation.
func (d *DNS) Validate() error {
	for _, z := range d.Zones {
		if err := z.Validate(); err != nil {
			return errors.Wrapf(err, "zone %s", z.Name)
		}
	}

	return nil
}

// Zone is a DNS zone.
type Zone struct {
	Name    string    `json:"name"`
	Records []*Record `json:"records"`
}

// Default implementation.
func (z *Zone) Default() error {
	for i, r := range z.Records {
		if err := r.Default(); err != nil {
			return errors.Wrapf(err, "record %d", i)
		}
	}

	return nil
}

// Validate implementation.
func (z *Zone) Validate() error {
	for i, r := range z.Records {
		if err := r.Validate(); err != nil {
			return errors.Wrapf(err, "record %d", i)
		}
	}

	return nil
}

// Record is a DNS record.
type Record struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	TTL   int      `json:"ttl"`
	Value []string `json:"value"`
}

// Validate implementation.
func (r *Record) Validate() error {
	if err := validate.List(r.Type, recordTypes); err != nil {
		return errors.Wrap(err, ".type")
	}

	if err := validate.RequiredString(r.Name); err != nil {
		return errors.Wrap(err, ".name")
	}

	if err := validate.RequiredStrings(r.Value); err != nil {
		return errors.Wrap(err, ".value")
	}

	if err := validate.MinStrings(r.Value, 1); err != nil {
		return errors.Wrap(err, ".value")
	}

	return nil
}

// Default implementation.
func (r *Record) Default() error {
	if r.TTL == 0 {
		r.TTL = 300
	}

	return nil
}
