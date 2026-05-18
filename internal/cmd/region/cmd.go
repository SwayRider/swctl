package region

import (
	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
)

var Command = &cli.Command{
	Name:  "region",
	Usage: "Query geographic regions via the RegionService",
	Flags: []cli.Flag{
		flags.Required(flags.Host("Region")),
		flags.Required(flags.Port("Region")),
		flags.Required(flags.Host("Auth")),
		flags.Required(flags.Port("Auth")),
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
	},
	Commands: []*cli.Command{
		SearchPoint,
		SearchBox,
		FindRegionPath,
	},
}
