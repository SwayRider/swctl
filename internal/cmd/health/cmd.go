package health

import (
	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
)

var Command = &cli.Command{
	Name:  "health",
	Usage: "Check the health of a service",
	Flags: []cli.Flag{
		flags.Required(flags.Host("Health")),
		flags.Required(flags.Port("Health")),
	},
	Commands: []*cli.Command{
		Ping,
		Check,
	},
}
