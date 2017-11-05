package app

import (
	"os"
	"time"

	"github.com/apex/up/internal/cli/root"

	"github.com/apex/up/internal/stats"
)

// Run the command.
func Run(version string) error {
	defer stats.Client.ConditionalFlush(100, 12*time.Hour)
	root.Cmd.Version(version)
	_, err := root.Cmd.Parse(os.Args[1:])
	return err
}
