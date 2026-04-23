package health

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/logic"
)

var Check = &cli.Command{
	Name:  "check",
	Usage: "Check the health status of a service or a specific component",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "component",
			Aliases: []string{"c"},
			Usage:   "The component to check (empty = overall service health)",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		status, err := logic.HealthCheck(
			c.String("health-host"),
			c.Int("health-port"),
			c.String("component"),
		)
		if err != nil {
			return err
		}
		fmt.Printf("status: %s\n", status)
		return nil
	},
}
