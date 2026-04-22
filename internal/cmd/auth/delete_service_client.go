package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var DeleteServiceClient = &cli.Command{
	Name:    "delete-service-client",
	Aliases: []string{"dsc"},
	Usage:   "Delete a service client",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "clientId",
			UsageText: "<clientId> The ClientID of the service client",
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		clientId := c.StringArg("clientId")
		if clientId == "" {
			return fmt.Errorf("clientId is required")
		}

		err := logic.DeleteServiceClient(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
			c.StringArg("clientId"),
		)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted service client %s\n", c.String("clientId"))
		return nil
	},
}
