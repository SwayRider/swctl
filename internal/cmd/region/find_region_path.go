package region

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/logic"
)

var FindRegionPath = &cli.Command{
	Name:    "find-region-path",
	Aliases: []string{"frp"},
	Usage:   "Find the region path between two regions",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "fromRegion",
			UsageText: "<fromRegion> The starting region code",
		},
		&cli.StringArg{
			Name:      "toRegion",
			UsageText: "<toRegion> The destination region code",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		from := c.StringArg("fromRegion")
		to := c.StringArg("toRegion")
		if from == "" || to == "" {
			return fmt.Errorf("fromRegion and toRegion are required")
		}

		path, err := logic.RegionFindPath(
			c.String("region-host"),
			c.Int("region-port"),
			from,
			to,
		)
		if err != nil {
			return err
		}

		if len(path) == 0 {
			fmt.Println("No path found.")
			return nil
		}
		fmt.Println(strings.Join(path, " → "))
		return nil
	},
}
