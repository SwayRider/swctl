package health

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/logic"
)

var Ping = &cli.Command{
	Name:  "ping",
	Usage: "Ping a service to check if it is reachable",
	Action: func(ctx context.Context, c *cli.Command) error {
		err := logic.Ping(
			c.String("health-host"),
			c.Int("health-port"),
		)
		if err != nil {
			return err
		}
		fmt.Println("pong")
		return nil
	},
}
