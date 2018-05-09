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

		assert.NoError(t, s.Default(), "default")
		assert.NoError(t, s.Validate(), "validate")

		assert.Len(t, s, 3)
		assert.Equal(t, "staging", s["staging"].Name)
		assert.Equal(t, "production", s["production"].Name)
	})

	t.Run("custom stages", func(t *testing.T) {
		s := Stages{
			"beta": &Stage{},
		}

		assert.NoError(t, s.Default(), "default")
		assert.NoError(t, s.Validate(), "validate")

		assert.Len(t, s, 4)
		assert.Equal(t, "beta", s["beta"].Name)
		assert.Equal(t, true, s["beta"].Zone)
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

		assert.NoError(t, s.Default(), "default")
		assert.NoError(t, s.Validate(), "validate")
		assert.Equal(t, "staging", s["staging"].Name)
		assert.Equal(t, "production", s["production"].Name)
	})

	t.Run("valid zone boolean", func(t *testing.T) {
		s := Stages{
			"production": &Stage{
				Domain: "gh-polls.com",
				Zone:   false,
			},
		}

		assert.NoError(t, s.Default(), "default")
		assert.NoError(t, s.Validate(), "validate")
	})

	t.Run("valid zone string", func(t *testing.T) {
		s := Stages{
			"production": &Stage{
				Domain: "api.gh-polls.com",
				Zone:   "api.gh-polls.com",
			},
		}

		assert.NoError(t, s.Default(), "default")
		assert.NoError(t, s.Validate(), "validate")
	})

	t.Run("invalid zone type", func(t *testing.T) {
		s := Stages{
			"production": &Stage{
				Domain: "api.gh-polls.com",
				Zone:   123,
			},
		}

		assert.NoError(t, s.Default(), "default")
		assert.EqualError(t, s.Validate(), `stage "production": .zone is an invalid type, must be string or boolean`)
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
