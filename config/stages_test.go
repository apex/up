package config

import (
	"testing"

	"github.com/tj/assert"
)

func TestStages_Validate(t *testing.T) {
	t.Run("no stages", func(t *testing.T) {
		s := Stages{}
		assert.NoError(t, s.Validate(), "validate")
	})

	t.Run("missing domain", func(t *testing.T) {
		s := Stages{
			Development: &Stage{},
		}

		assert.EqualError(t, s.Validate(), `.development: .domain: is required`)
	})

	t.Run("some stages", func(t *testing.T) {
		s := Stages{
			Development: &Stage{
				Domain: "gh-polls-dev.com",
			},
			Staging: &Stage{
				Domain: "gh-polls-stage.com",
			},
			Production: &Stage{
				Domain: "gh-polls.com",
			},
		}

		assert.NoError(t, s.Default(), "validate")
		assert.NoError(t, s.Validate(), "validate")
		assert.Equal(t, "development", s.Development.Name)
		assert.Equal(t, "staging", s.Staging.Name)
		assert.Equal(t, "production", s.Production.Name)
	})
}

func TestStages_List(t *testing.T) {
	dev := &Stage{
		Domain: "gh-polls-dev.com",
	}

	prod := &Stage{
		Domain: "gh-polls.com",
	}

	s := Stages{
		Development: dev,
		Production:  prod,
	}

	list := []*Stage{
		dev,
		prod,
	}

	stages := s.List()
	assert.Equal(t, list, stages)
}

func TestStages_GetByDomain(t *testing.T) {
	dev := &Stage{
		Domain: "gh-polls-dev.com",
	}

	prod := &Stage{
		Domain: "gh-polls.com",
	}

	s := Stages{
		Development: dev,
		Production:  prod,
	}

	stage := s.GetByDomain("gh-polls.com")
	assert.Equal(t, prod, stage)
}
