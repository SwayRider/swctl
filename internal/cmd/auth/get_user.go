package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var GetUser = &cli.Command{
	Name:    "get-user",
	Aliases: []string{"gu"},
	Usage:   "Get user info by email or user ID",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "identifier",
			UsageText: "<email|userId> The email address or user ID of the user",
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		identifier := c.StringArg("identifier")
		if identifier == "" {
			return fmt.Errorf("identifier (email or user ID) is required")
		}

		user, err := logic.GetUser(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
			identifier,
		)
		if err != nil {
			return err
		}
		user.Display()
		return nil
	},
}
