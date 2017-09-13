package up

import (
	"io/ioutil"
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

func TestConfig_Endpoint(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		endpoint := "http://localhost:1234"
		c := Config{
			Name:     "my-app123",
			Endpoint: endpoint,
		}

		assert.NoError(t, c.Default(), "default")
		assert.NoError(t, c.Validate(), "validate")
		assert.Equal(t, endpoint, c.Endpoint)
	})

	t.Run("invalid", func(t *testing.T) {
		c := Config{
			Name:     "my-app123",
			Endpoint: "http//localhost",
		}

		assert.NoError(t, c.Default(), "default")
		assert.EqualError(t, c.Validate(), `.endpoint: url: parse http//localhost: invalid URI for request`)
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
  • eu-west-2
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
  • eu-west-2
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
		assert.NoError(t, c.Default(), "default")

		assert.NoError(t, c.defaultRegions(), "defaultRegions")
		assert.Equal(t, 2, len(c.Regions), "regions should have length 2")
		assert.Equal(t, regions, c.Regions, "should read regions from config")
		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("regions from AWS_REGION", func(t *testing.T) {
		region := "sa-east-1"
		os.Setenv("AWS_REGION", region)

		defer os.Setenv("AWS_REGION", "")
		c := Config{
			Name: "api",
			Type: "server",
		}
		assert.NoError(t, c.Default(), "default")

		assert.NoError(t, c.defaultRegions(), "defaultRegions")
		assert.Equal(t, 1, len(c.Regions), "regions should have length 1")
		assert.Equal(t, region, c.Regions[0], "should read regions from AWS_REGION")
		assert.NoError(t, c.Default(), "default")
		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("regions from AWS_DEFAULT_REGION", func(t *testing.T) {
		region := "sa-east-1"

		os.Setenv("AWS_DEFAULT_REGION", region)
		defer os.Setenv("AWS_DEFAULT_REGION", "")

		c := Config{
			Name: "api",
			Type: "server",
		}
		assert.NoError(t, c.Default(), "default")

		assert.NoError(t, c.defaultRegions(), "defaultRegions")
		assert.Equal(t, 1, len(c.Regions), "regions should have length 1")
		assert.Equal(t, region, c.Regions[0], "should read regions from AWS_DEFAULT_REGION")
		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("regions from shared config with default profile", func(t *testing.T) {
		content := `
		[default]
		region = sa-east-1
		output = json
		[profile another-profile]
		region = ap-southeast-2
		output = json`

		tmpfile, err := ioutil.TempFile("", "config")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		_, err = tmpfile.WriteString(content)
		assert.NoError(t, err)

		os.Setenv("AWS_CONFIG_FILE", tmpfile.Name())
		defer os.Setenv("AWS_CONFIG_FILE", "")

		c := Config{
			Name: "api",
			Type: "server",
		}
		assert.NoError(t, c.Default(), "default")

		assert.NoError(t, c.defaultRegions(), "defaultRegions")
		assert.Equal(t, 1, len(c.Regions), "regions should have length 1")
		assert.Equal(t, "sa-east-1", c.Regions[0], "should read regions from shared config with default profile")
		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("regions from shared config with AWS_PROFILE profile", func(t *testing.T) {
		content := `
		[default]
		region = sa-east-1
		output = json
		[profile another-profile]
		region = ap-southeast-2
		output = json`

		tmpfile, err := ioutil.TempFile("", "config")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		_, err = tmpfile.WriteString(content)
		assert.NoError(t, err)

		os.Setenv("AWS_CONFIG_FILE", tmpfile.Name())
		defer os.Setenv("AWS_CONFIG_FILE", "")

		os.Setenv("AWS_PROFILE", "another-profile")
		defer os.Setenv("AWS_PROFILE", "")

		c := Config{
			Name: "api",
			Type: "server",
		}
		assert.NoError(t, c.Default(), "default")

		assert.NoError(t, c.defaultRegions(), "defaultRegions")
		assert.Equal(t, 1, len(c.Regions), "regions should have length 1")
		assert.Equal(t, "ap-southeast-2", c.Regions[0], "should read regions from shared config with AWS_PROFILE profile")
		assert.NoError(t, c.Validate(), "validate")
	})

	t.Run("default region must be us-west-2", func(t *testing.T) {
		// Make sure we aren't reading AWS config file
		os.Setenv("AWS_CONFIG_FILE", "does-not-exist")

		c := Config{
			Name: "api",
			Type: "server",
		}
		assert.NoError(t, c.Default(), "default")

		assert.NoError(t, c.defaultRegions(), "defaultRegions")
		assert.Equal(t, 1, len(c.Regions), "regions should have length 1")
		assert.Equal(t, "us-west-2", c.Regions[0], "default region must be us-west-2")
		assert.NoError(t, c.Validate(), "validate")
	})
}
