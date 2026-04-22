package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
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
			UsageText: "<password> The password of the admin",
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
		pwd := c.StringArg("password")
		if pwd == "" {
			return fmt.Errorf("password is required")
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
