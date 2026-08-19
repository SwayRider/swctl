package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
	"github.com/swayrider/swctl/internal/prompt"
)

var CreateAdmin = &cli.Command{
	Name:    "create-admin",
	Aliases: []string{"ca"},
	Usage:   "Crate a new admin user",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "email",
			UsageText: "<email> The email of the admin ",
		},
		&cli.StringArg{
			Name:      "password",
			UsageText: "<password> The password of the admin (optional; prompted if omitted)",
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
	},
	Before: prompt.BeforeFillPassword,
	Action: func(ctx context.Context, c *cli.Command) error {
		email := c.StringArg("email")
		if email == "" {
			return fmt.Errorf("email is required")
		}
		pwd := c.StringArg("password")
		if pwd == "" {
			var err error
			pwd, err = prompt.Password("Password for new admin")
			if err != nil {
				return err
			}
		}

		user, err := logic.CreateAdmin(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
			email,
			pwd,
		)
		if err != nil {
			return err
		}
		fmt.Println("Created admin:")
		user.Display()
		return nil
	},
}
