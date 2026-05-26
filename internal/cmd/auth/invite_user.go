package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var InviteUser = &cli.Command{
	Name:    "invite-user",
	Aliases: []string{"iu"},
	Usage:   "Add an email address to the registration invite list",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "email",
			UsageText: "<email> The email address to invite",
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

		msg, err := logic.InviteUser(
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
