package config

import (
	"testing"

	"github.com/tj/assert"
)

func TestStage_Override(t *testing.T) {
	c, err := ParseConfigString(`{
		"name": "app",
		"regions": ["us-west-2"],
		"lambda": {
			"memory": 128
		},
		"hooks": {
			"build": "parcel index.html -o build",
			"clean": "rm -fr build"
		},
		"proxy": {
			"command": "node app.js"
		},
		"stages": {
			"production": {
				"lambda": {
					"memory": 1024
				},
				"hooks": {
					"build": "parcel index.html -o build --production"
				}
			},
			"staging": {
				"proxy": {
					"command": "node app.js --foo=bar"
				}
			}
		}
	}`)

	assert.NoError(t, err, "parse")
	assert.NoError(t, c.Default(), "default")
	assert.NoError(t, c.Validate(), "validate")

	assert.NoError(t, c.Override("production"), "override")
	assert.Equal(t, 1024, c.Lambda.Memory)
	assert.Equal(t, Hook{`parcel index.html -o build --production`}, c.Hooks.Build)
	assert.Equal(t, `node app.js`, c.Proxy.Command)

	assert.NoError(t, c.Override("staging"), "override")
	assert.Equal(t, `node app.js --foo=bar`, c.Proxy.Command)
}

func TestStages_Default(t *testing.T) {
	t.Run("no custom stages", func(t *testing.T) {
		s := Stages{}

		assert.NoError(t, s.Default(), "validate")
		assert.NoError(t, s.Validate(), "validate")

		assert.Len(t, s, 3)
		assert.Equal(t, "staging", s["staging"].Name)
		assert.Equal(t, "production", s["production"].Name)
	})

	t.Run("custom stages", func(t *testing.T) {
		s := Stages{
			"beta": &Stage{},
		}

		assert.NoError(t, s.Default(), "validate")
		assert.NoError(t, s.Validate(), "validate")

		assert.Len(t, s, 4)
		assert.Equal(t, "beta", s["beta"].Name)
	})
}

func TestStages_Validate(t *testing.T) {
	t.Run("no stages", func(t *testing.T) {
		s := Stages{}
		assert.NoError(t, s.Validate(), "validate")
	})

	t.Run("some stages", func(t *testing.T) {
		s := Stages{
			"staging": &Stage{
				Domain: "gh-polls-stage.com",
			},
			"production": &Stage{
				Domain: "gh-polls.com",
			},
		}

		assert.NoError(t, s.Default(), "validate")
		assert.NoError(t, s.Validate(), "validate")
		assert.Equal(t, "staging", s["staging"].Name)
		assert.Equal(t, "production", s["production"].Name)
	})
}

func TestStages_List(t *testing.T) {
	stage := &Stage{
		Domain: "gh-polls-stage.com",
	}

	prod := &Stage{
		Domain: "gh-polls.com",
	}

	s := Stages{
		"staging":    stage,
		"production": prod,
	}

	list := []*Stage{
		stage,
		prod,
	}

	stages := s.List()
	assert.Equal(t, list, stages)
}

func TestStages_GetByDomain(t *testing.T) {
	stage := &Stage{
		Domain: "gh-polls-stage.com",
	}

	prod := &Stage{
		Domain: "gh-polls.com",
	}

	s := Stages{
		"staging":    stage,
		"production": prod,
	}

	assert.Equal(t, prod, s.GetByDomain("gh-polls.com"))
}
