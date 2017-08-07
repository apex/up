package text

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/apex/log"
)

func init() {
	log.Now = func() time.Time {
		return time.Unix(0, 0)
	}
}

func Test(t *testing.T) {
	var buf bytes.Buffer

	log.SetHandler(New(&buf))
	log.WithField("user", "tj").WithField("id", "123").Info("hello")
	log.WithField("user", "tj").Info("something broke")
	log.WithField("user", "tj").Warn("something kind of broke")
	log.WithField("user", "tj").Error("boom")

	io.Copy(os.Stdout, &buf)
}
