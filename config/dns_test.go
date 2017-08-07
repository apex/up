package config

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/tj/assert"
)

func ExampleDNS() {
	s := `{
		"something.sh": [
			{
				"name": "something.com",
				"type": "A",
				"ttl": 300,
				"value": ["35.161.83.243"]
			},
			{
				"name": "blog.something.com",
				"type": "CNAME",
				"ttl": 300,
				"value": ["34.209.172.67"]
			},
			{
				"name": "api.something.com",
				"type": "A",
				"ttl": 300,
				"value": ["54.187.185.18"]
			}
		]
	}`

	var c DNS

	if err := json.Unmarshal([]byte(s), &c); err != nil {
		log.Fatalf("error unmarshaling: %s", err)
	}

	sort.Slice(c.Zones[0].Records, func(i int, j int) bool {
		a := c.Zones[0].Records[i]
		b := c.Zones[0].Records[j]
		return a.Name > b.Name
	})

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(c)
	// Output:
	// 	{
	//   "zones": [
	//     {
	//       "name": "something.sh",
	//       "records": [
	//         {
	//           "name": "something.com",
	//           "type": "A",
	//           "ttl": 300,
	//           "value": [
	//             "35.161.83.243"
	//           ]
	//         },
	//         {
	//           "name": "blog.something.com",
	//           "type": "CNAME",
	//           "ttl": 300,
	//           "value": [
	//             "34.209.172.67"
	//           ]
	//         },
	//         {
	//           "name": "api.something.com",
	//           "type": "A",
	//           "ttl": 300,
	//           "value": [
	//             "54.187.185.18"
	//           ]
	//         }
	//       ]
	//     }
	//   ]
	// }
}

func TestDNS_Validate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		c := &DNS{
			Zones: []*Zone{
				{
					Name: "apex.sh",
					Records: []Record{
						{
							Name: "blog.apex.sh",
							Type: "CNAME",
						},
					},
				},
			},
		}

		assert.EqualError(t, c.Validate(), `zone apex.sh: record 0: .value: must have at least 1 value`)
	})
}

func TestRecord_Type(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		c := &Record{
			Name:  "blog.apex.sh",
			Type:  "A",
			Value: []string{"1.1.1.1"},
		}

		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("invalid", func(t *testing.T) {
		c := &Record{
			Name: "blog.apex.sh",
			Type: "AAA",
		}

		assert.EqualError(t, c.Validate(), `.type: "AAA" is invalid, must be one of:

  • ALIAS
  • A
  • AAAA
  • CNAME
  • MX
  • NAPTR
  • NS
  • PTR
  • SOA
  • SPF
  • SRV
  • TXT`)
	})
}

func TestRecord_TTL(t *testing.T) {
	c := &Record{Type: "A"}
	assert.NoError(t, c.Default(), "default")
	assert.Equal(t, 300, c.TTL)
}

func TestRecord_Value(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		c := &Record{
			Name: "blog.apex.sh",
			Type: "A",
		}

		assert.EqualError(t, c.Validate(), `.value: must have at least 1 value`)
	})

	t.Run("invalid", func(t *testing.T) {
		c := &Record{
			Name:  "blog.apex.sh",
			Type:  "A",
			Value: []string{"1.1.1.1", ""},
		}

		assert.EqualError(t, c.Validate(), `.value: at index 1: is required`)
	})
}
