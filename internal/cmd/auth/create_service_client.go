package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var CreateServiceClient = &cli.Command{
	Name:    "create-service-client",
	Aliases: []string{"csc"},
	Usage:   "Create a service client",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "name",
			UsageText: "<name> The name of the service client ",
		},
		&cli.StringArgs{
			Name:      "scope",
			UsageText: "<scope...> The scopes for the service client",
			Min:       1,
			Max:       -1,
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"d"},
			Usage:   "The description of the service client",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		name := c.StringArg("name")
		if name == "" {
			return fmt.Errorf("name is required")
		}

		clientId, clientSecret, err := logic.CreateServiceClient(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
			name,
			c.String("description"),
			c.StringArgs("scope"),
		)
		if err != nil {
			return err
		}

		fmt.Println("WARNING: Store the client secret securely — it will not be shown again.")
		fmt.Printf("client id: %s\nclient secret: %s\n", clientId, clientSecret)
		return nil
	},
}
