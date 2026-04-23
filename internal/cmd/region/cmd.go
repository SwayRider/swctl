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
	},
	Commands: []*cli.Command{
		SearchPoint,
		SearchBox,
		FindRegionPath,
	},
}
