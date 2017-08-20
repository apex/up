package up

import (
	"os"
	"testing"

	"github.com/tj/assert"
)

func TestConfig_Name(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		c := Config{
			Name: "my-app123",
		}

		assert.NoError(t, c.Default(), "default")
		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("invalid", func(t *testing.T) {
		c := Config{
			Name: "my app",
		}

		assert.NoError(t, c.Default(), "default")
		assert.EqualError(t, c.Validate(), `.name "my app": must contain only alphanumeric characters and '-'`)
	})
}

func TestConfig_Type(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		c := Config{
			Name: "api",
		}

		assert.NoError(t, c.Default(), "default")
		assert.NoError(t, c.Validate(), "validate")
		assert.Equal(t, "server", c.Type)
	})

	t.Run("valid", func(t *testing.T) {
		c := Config{
			Name: "api",
			Type: "server",
		}

		assert.NoError(t, c.Default(), "default")
		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("invalid", func(t *testing.T) {
		c := Config{
			Name: "api",
			Type: "something",
		}

		assert.NoError(t, c.Default(), "default")
		assert.EqualError(t, c.Validate(), `.type: "something" is invalid, must be one of:

  • static
  • server`)
	})
}

func TestConfig_Regions(t *testing.T) {
	t.Run("valid multiple", func(t *testing.T) {
		c := Config{
			Name:    "api",
			Type:    "server",
			Regions: []string{"us-west-2", "us-east-1"},
		}

		assert.NoError(t, c.Default(), "default")
		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("valid multiple", func(t *testing.T) {
		c := Config{
			Name:    "api",
			Type:    "server",
			Regions: []string{"us-west-2", "us-east-1"},
		}

		assert.NoError(t, c.Default(), "default")
		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("valid globbing", func(t *testing.T) {
		c := Config{
			Name:    "api",
			Type:    "server",
			Regions: []string{"us-*", "us-east-1", "ca-central-*"},
		}

		assert.NoError(t, c.Default(), "default")
		assert.NoError(t, c.Validate(), "validate")
		assert.Equal(t, []string{"us-east-1", "us-west-2", "us-east-2", "us-west-1", "us-east-1", "ca-central-1"}, c.Regions)
	})

	t.Run("invalid globbing", func(t *testing.T) {
		c := Config{
			Name:    "api",
			Type:    "server",
			Regions: []string{"uss-*"},
		}

		assert.NoError(t, c.Default(), "default")

		assert.EqualError(t, c.Validate(), `.regions: "uss-*" is invalid, must be one of:

  • us-east-1
  • us-west-2
  • eu-west-1
  • eu-central-1
  • ap-northeast-1
  • ap-southeast-1
  • ap-southeast-2
  • us-east-2
  • us-west-1
  • ap-northeast-2
  • ap-south-1
  • sa-east-1
  • ca-central-1`)
	})

	t.Run("invalid", func(t *testing.T) {
		c := Config{
			Name:    "api",
			Type:    "server",
			Regions: []string{"us-west-1", "us-west-9"},
		}

		assert.NoError(t, c.Default(), "default")

		assert.EqualError(t, c.Validate(), `.regions: "us-west-9" is invalid, must be one of:

  • us-east-1
  • us-west-2
  • eu-west-1
  • eu-central-1
  • ap-northeast-1
  • ap-southeast-1
  • ap-southeast-2
  • us-east-2
  • us-west-1
  • ap-northeast-2
  • ap-south-1
  • sa-east-1
  • ca-central-1`)
	})
}

func TestConfig_defaultRegions(t *testing.T) {
	t.Run("regions from config", func(t *testing.T) {
		regions := []string{"us-west-2", "us-east-1"}
		c := Config{
			Name:    "api",
			Type:    "server",
			Regions: regions,
		}

		assert.NoError(t, c.defaultRegions(), "defaultRegions")
		assert.Equal(t, regions, c.Regions, "should read regions from config")
	})

	t.Run("regions from AWS_REGION", func(t *testing.T) {
		region := "sa-east-1"
		os.Setenv("AWS_REGION", region)
		defer os.Setenv("AWS_REGION", "")
		c := Config{
			Name: "api",
			Type: "server",
		}

		assert.NoError(t, c.defaultRegions(), "defaultRegions")
		assert.Equal(t, 1, len(c.Regions), "regions should have length 1")
		assert.Equal(t, region, c.Regions[0], "should read regions from AWS_REGION")
	})

	t.Run("regions from AWS_DEFAULT_REGION", func(t *testing.T) {
		region := "sa-east-1"
		os.Setenv("AWS_DEFAULT_REGION", region)
		defer os.Setenv("AWS_DEFAULT_REGION", "")
		c := Config{
			Name: "api",
			Type: "server",
		}

		assert.NoError(t, c.defaultRegions(), "defaultRegions")
		assert.Equal(t, 1, len(c.Regions), "regions should have length 1")
		assert.Equal(t, region, c.Regions[0], "should read regions from AWS_DEFAULT_REGION")
	})

}
