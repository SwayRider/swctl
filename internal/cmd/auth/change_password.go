package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
	"github.com/swayrider/swctl/internal/prompt"
)

var ChangePassword = &cli.Command{
	Name:    "change-password",
	Aliases: []string{"chp"},
	Usage:   "Change the password for the authenticated user (--password is used as the current/old password)",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "newPassword",
			UsageText: "<newPassword> The new password (optional; prompted if omitted)",
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
	},
	Before: prompt.BeforeFillPassword,
	Action: func(ctx context.Context, c *cli.Command) error {
		newPassword := c.StringArg("newPassword")
		if newPassword == "" {
			var err error
			newPassword, err = prompt.Password("New password")
			if err != nil {
				return err
			}
		}

		msg, err := logic.ChangePassword(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
			newPassword,
		)
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	},
}
