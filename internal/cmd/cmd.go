package cmd

import (
	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/cmd/auth"
)

var App = &cli.Command{
	Name:  "swctl",
	Usage: "SwayRider service control tool",
	Flags: []cli.Flag{},
	Commands: []*cli.Command{
		auth.Command,
	},
}
