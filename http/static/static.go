package static

import (
	"net/http"

	"github.com/apex/up"
)

// New static handler.
func New(c *up.Config) http.Handler {
	return http.FileServer(http.Dir(c.Static.Dir))
}
