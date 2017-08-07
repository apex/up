package logs

import (
	"io"
	"os"
	"time"

	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/stats"
	"github.com/tj/kingpin"
)

func init() {
	cmd := root.Command("logs", "Show log output.")
	cmd.Example(`up logs`, "Show logs from the past 5 minutes.")
	cmd.Example(`up logs -s 30m`, "Show logs from the past 30 minutes.")
	cmd.Example(`up logs -s 5h`, "Show logs from the past 5 hours.")
	cmd.Example(`up logs -f`, "Show live log output.")
	cmd.Example(`up logs error`, "Show error logs.")
	cmd.Example(`up logs 'error or fatal'`, "Show error and fatal logs.")
	cmd.Example(`up logs 'not info'`, "Show non-info logs.")
	cmd.Example(`up logs 'message = "user login"'`, "Show logs with a specific message.")
	cmd.Example(`up logs 'status = 200 duration > 150'`, "Show 200 responses with latency above 150ms.")
	cmd.Example(`up logs 'status >= 400'`, "Show 4xx and 5xx responses.")
	cmd.Example(`up logs 'user.email contains "@apex.sh"'`, "Show emails containing @apex.sh.")
	cmd.Example(`up logs 'user.email = "*@apex.sh"'`, "Show emails ending with @apex.sh.")
	cmd.Example(`up logs 'user.email = "tj@*"'`, "Show emails starting with tj@.")
	cmd.Example(`up logs 'method in ("POST", "PUT") ip = "207.*" status = 200 duration >= 50'`, "Show logs with a more complex query.")
	cmd.Example(`up logs error | jq`, "Pipe JSON error logs to the jq tool.")

	// TODO: query arg should join remaining args
	query := cmd.Arg("query", "Query pattern for filtering logs.").String()
	follow := cmd.Flag("follow", "Follow or tail the live logs.").Short('f').Bool()
	since := cmd.Flag("since", "Show logs since duration (30s, 5m, 2h, 1h30m).").Short('s').Default("5m").Duration()

	cmd.Action(func(_ *kingpin.ParseContext) error {
		q := *query
		s := *since

		if *follow {
			s = time.Duration(0)
		}

		stats.Track("Logs", map[string]interface{}{
			"query":        q != "",
			"query_length": len(q),
			"follow":       *follow,
			"since":        s.Round(time.Second),
		})

		logs := root.Project.Logs(q)
		logs.Since(time.Now().Add(-s))

		if *follow {
			logs.Follow()
		}

		_, err := io.Copy(os.Stdout, logs)
		return err
	})
}
