package logs

import (
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/tj/go/term"
	"github.com/tj/kingpin"

	"github.com/apex/up"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/util"
)

func init() {
	cmd := root.Command("logs", "Show log output.")
	cmd.Example(`up logs`, "Show logs from the past hour.")
	cmd.Example(`up logs -s 30m`, "Show logs from the past 30 minutes.")
	cmd.Example(`up logs -s 5h`, "Show logs from the past 5 hours.")
	cmd.Example(`up logs -f`, "Show live log output.")
	cmd.Example(`up logs error`, "Show error logs.")
	cmd.Example(`up logs 'production (warn or error)'`, "Show 4xx and 5xx responses in production.")
	cmd.Example(`up logs 'production error method in ("POST", "PUT", "DELETE")'`, "Show production 5xx responses with a POST, PUT, or DELETE method.")
	cmd.Example(`up logs 'error or fatal'`, "Show error and fatal logs.")
	cmd.Example(`up logs 'not info'`, "Show non-info logs.")
	cmd.Example(`up logs 'message = "user login"'`, "Show logs with a specific message.")
	cmd.Example(`up logs 'status = 200 duration > 1.5s'`, "Show 200 responses with latency above 1500ms.")
	cmd.Example(`up logs 'size > 100kb'`, "Show responses with bodies larger than 100kb.")
	cmd.Example(`up logs 'status >= 400'`, "Show 4xx and 5xx responses.")
	cmd.Example(`up logs 'user.email contains "@apex.sh"'`, "Show emails containing @apex.sh.")
	cmd.Example(`up logs 'user.email = "*@apex.sh"'`, "Show emails ending with @apex.sh.")
	cmd.Example(`up logs 'user.email = "tj@*"'`, "Show emails starting with tj@.")
	cmd.Example(`up logs 'method in ("POST", "PUT") ip = "207.*" status = 200 duration >= 50'`, "Show logs with a more complex query.")
	cmd.Example(`up logs error | jq`, "Pipe JSON error logs to the jq tool.")

	query := cmd.Arg("query", "Query pattern for filtering logs.").String()
	follow := cmd.Flag("follow", "Follow or tail the live logs.").Short('f').Bool()
	since := cmd.Flag("since", "Show logs since duration (30s, 5m, 2h, 1h30m, 3d, 1M).").Short('s').Default("1d").String()
	expand := cmd.Flag("expand", "Show expanded logs.").Short('e').Bool()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		var s time.Duration

		if *since != "" {
			s, err = util.ParseDuration(*since)
			if err != nil {
				return errors.Wrap(err, "parsing --since duration")
			}
		}

		if *follow {
			s = time.Duration(0)
		}

		q := *query

		stats.Track("Logs", map[string]interface{}{
			"query":        q != "",
			"query_length": len(q),
			"follow":       *follow,
			"since":        s.Round(time.Second),
			"expand":       *expand,
		})

		logs := p.Logs(up.LogsConfig{
			Region:     c.Regions[0],
			Since:      time.Now().Add(-s),
			Follow:     *follow,
			Expand:     *expand,
			Query:      q,
			OutputJSON: !term.IsTerminal(os.Stdout.Fd()),
		})

		if _, err := io.Copy(os.Stdout, logs); err != nil {
			return errors.Wrap(err, "writing logs")
		}

		return nil
	})
}
