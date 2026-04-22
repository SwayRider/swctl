package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var CreateUser = &cli.Command{
	Name:    "create-user",
	Aliases: []string{"cu"},
	Usage:   "Create a new user with an option account type (default: free) and optionally set it as verified",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "email",
			UsageText: "<email> The email of the user ",
		},
		&cli.StringArg{
			Name:      "password",
			UsageText: "<password> The password of the user",
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
		&cli.BoolFlag{
			Name:     "verified",
			Aliases:  []string{"v"},
			Usage:    "Set the user as verified",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "account-type",
			Aliases:  []string{"t"},
			Usage:    "Set the account type (default: free)",
			Required: false,
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		email := c.StringArg("email")
		if email == "" {
			return fmt.Errorf("email is required")
		}
		password := c.StringArg("password")
		if password == "" {
			return fmt.Errorf("password is required")
		}

		user, err := logic.CreateUser(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
			email,
			password,
			c.Bool("verified"),
			c.String("account-type"),
		)
		if err != nil {
			return err
		}
		fmt.Println("Created user:")
		user.Display()
		return nil
	},
}
