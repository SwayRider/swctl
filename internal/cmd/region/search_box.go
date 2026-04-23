package region

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/logic"
)

var SearchBox = &cli.Command{
	Name:    "search-box",
	Aliases: []string{"sb"},
	Usage:   "Find regions intersecting a bounding box",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "minLat",
			UsageText: "<minLat> Bottom-left latitude",
		},
		&cli.StringArg{
			Name:      "minLon",
			UsageText: "<minLon> Bottom-left longitude",
		},
		&cli.StringArg{
			Name:      "maxLat",
			UsageText: "<maxLat> Top-right latitude",
		},
		&cli.StringArg{
			Name:      "maxLon",
			UsageText: "<maxLon> Top-right longitude",
		},
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "extended",
			Aliases: []string{"e"},
			Usage:   "Include extended region data",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		parse := func(name, val string) (float64, error) {
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid %s: %w", name, err)
			}
			return v, nil
		}

		minLat, err := parse("minLat", c.StringArg("minLat"))
		if err != nil {
			return err
		}
		minLon, err := parse("minLon", c.StringArg("minLon"))
		if err != nil {
			return err
		}
		maxLat, err := parse("maxLat", c.StringArg("maxLat"))
		if err != nil {
			return err
		}
		maxLon, err := parse("maxLon", c.StringArg("maxLon"))
		if err != nil {
			return err
		}

		regions, err := logic.RegionSearchBox(
			c.String("region-host"),
			c.Int("region-port"),
			minLat, minLon, maxLat, maxLon,
			c.Bool("extended"),
		)
		if err != nil {
			return err
		}

		fmt.Printf("Core regions:     %s\n", strings.Join(regions.CoreRegions, ", "))
		if c.Bool("extended") {
			fmt.Printf("Extended regions: %s\n", strings.Join(regions.ExtendedRegions, ", "))
		}
		return nil
	},
}
