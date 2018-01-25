package proxy

import (
	"testing"

	"github.com/tj/assert"
)

func Test_JSON_isTextMime(t *testing.T) {
	assert.Equal(t, isTextMime("application/json"), true)
	assert.Equal(t, isTextMime("application/json; charset=utf-8"), true)
	assert.Equal(t, isTextMime("Application/JSON"), true)
}

func Test_XML_isTextMime(t *testing.T) {
	assert.Equal(t, isTextMime("application/xml"), true)
	assert.Equal(t, isTextMime("application/xml; charset=utf-8"), true)
	assert.Equal(t, isTextMime("ApPlicaTion/xMl"), true)
}
