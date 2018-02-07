package config

import (
	"testing"

	"github.com/tj/assert"
)

func TestStage_Override(t *testing.T) {
	c := &Config{
		Regions: []string{"us-west-2"},
		Lambda: Lambda{
			Memory: 128,
			Role:   "arn::something",
		},
		Hooks: Hooks{
			Build: Hook{"parcel index.html -o build"},
			Clean: Hook{"rm -fr clean"},
		},
		Proxy: Relay{
			Command: "./server",
		},
	}

	t.Run("with no overrides", func(t *testing.T) {
		c := *c

		s := &Stage{
			Domain: "example.com",
		}

		s.Override(&c)

		assert.Equal(t, []string{"us-west-2"}, c.Regions)
		assert.Equal(t, 128, c.Lambda.Memory)
		assert.Equal(t, "arn::something", c.Lambda.Role)

		assert.Equal(t, "parcel index.html -o build", c.Hooks.Build[0])
		assert.Equal(t, "rm -fr clean", c.Hooks.Clean[0])
		assert.Equal(t, `./server`, c.Proxy.Command)
	})

	t.Run("with overrides", func(t *testing.T) {
		c := *c

		s := &Stage{
			Domain: "example.com",
			StageOverrides: StageOverrides{
				Hooks: Hooks{
					Build:     Hook{"parcel index.html -o build --production"},
					PostBuild: Hook{"do something"},
				},
				Lambda: Lambda{
					Memory: 1024,
				},
				Proxy: Relay{
					Command: "./server --foo",
				},
			},
		}

		s.Override(&c)

		assert.Equal(t, 1024, c.Lambda.Memory)
		assert.Equal(t, "arn::something", c.Lambda.Role)

		assert.Equal(t, "parcel index.html -o build --production", c.Hooks.Build[0])
		assert.Equal(t, "do something", c.Hooks.PostBuild[0])
		assert.Equal(t, "rm -fr clean", c.Hooks.Clean[0])
		assert.Equal(t, `./server --foo`, c.Proxy.Command)
	})
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
