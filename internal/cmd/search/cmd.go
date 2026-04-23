package search

import (
	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
)

var Command = &cli.Command{
	Name:  "search",
	Usage: "Geocode and reverse geocode via the SearchService",
	Flags: []cli.Flag{
		flags.Required(flags.Host("Auth")),
		flags.Required(flags.Port("Auth")),
		flags.Required(flags.Host("Search")),
		flags.Required(flags.Port("Search")),
	},
	Commands: []*cli.Command{
		Geocode,
		Reverse,
	},
}
