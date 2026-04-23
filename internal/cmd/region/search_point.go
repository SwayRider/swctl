package region

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/logic"
)

var SearchPoint = &cli.Command{
	Name:    "search-point",
	Aliases: []string{"sp"},
	Usage:   "Find regions containing a coordinate",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "lat",
			UsageText: "<lat> Latitude",
		},
		&cli.StringArg{
			Name:      "lon",
			UsageText: "<lon> Longitude",
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
		latStr := c.StringArg("lat")
		lonStr := c.StringArg("lon")
		if latStr == "" || lonStr == "" {
			return fmt.Errorf("lat and lon are required")
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return fmt.Errorf("invalid lat: %w", err)
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return fmt.Errorf("invalid lon: %w", err)
		}

		regions, err := logic.RegionSearchPoint(
			c.String("region-host"),
			c.Int("region-port"),
			lat, lon,
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
