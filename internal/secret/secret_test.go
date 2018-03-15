package secret

import (
	"testing"

	"github.com/apex/up"
	"github.com/tj/assert"
)

func TestFormat(t *testing.T) {
	t.Run("relative", func(t *testing.T) {
		assert.Equal(t, "/up/myapp/production/S3_KEY", Format("myapp", "production", "S3_KEY"))
	})

	t.Run("absolute", func(t *testing.T) {
		assert.Equal(t, "/up/myapp/all/S3_KEY", Format("myapp", "", "S3_KEY"))
	})
}

func TestParse(t *testing.T) {
	s := Format("myapp", "production", "S3_KEY")
	app, stage, name := Parse(s)
	assert.Equal(t, "myapp", app)
	assert.Equal(t, "production", stage)
	assert.Equal(t, "S3_KEY", name)
}

func TestEnv(t *testing.T) {
	s := []*up.Secret{
		{
			Name:  "foo",
			Value: "bar",
		},
		{
			Name:  "bar",
			Value: "baz",
		},
	}

	env := Env(s)
	assert.Equal(t, []string{"foo=bar", "bar=baz"}, env)
}

func TestString(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		s := &up.Secret{
			Type:  "String",
			Value: "hello",
		}

		assert.Equal(t, `hello`, String(s))
	})

	t.Run("SecureString", func(t *testing.T) {
		s := &up.Secret{
			Type:  "SecureString",
			Value: "hello",
		}

		assert.Equal(t, `*****`, String(s))
	})
}
