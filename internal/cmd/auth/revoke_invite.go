package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var RevokeInvite = &cli.Command{
	Name:    "revoke-invite",
	Aliases: []string{"ri"},
	Usage:   "Remove an email address from the registration invite list",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "email",
			UsageText: "<email> The email address to remove from the invite list",
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		email := c.StringArg("email")
		if email == "" {
			return fmt.Errorf("email is required")
		}

		msg, err := logic.RevokeInvite(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
			email,
		)
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	},
}
