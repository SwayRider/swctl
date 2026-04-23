package cmd

import (
	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/cmd/auth"
	"github.com/swayrider/swctl/internal/cmd/health"
	"github.com/swayrider/swctl/internal/cmd/region"
	"github.com/swayrider/swctl/internal/cmd/search"
)

var App = &cli.Command{
	Name:  "swctl",
	Usage: "SwayRider service control tool",
	Flags: []cli.Flag{},
	Commands: []*cli.Command{
		auth.Command,
		health.Command,
		search.Command,
		region.Command,
	},
}
