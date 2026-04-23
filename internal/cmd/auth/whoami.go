package auth

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var WhoAmI = &cli.Command{
	Name:    "whoami",
	Aliases: []string{"wai"},
	Usage:   "Show the identity of the configured credentials",
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		user, err := logic.WhoAmI(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
		)
		if err != nil {
			return err
		}
		user.Display()
		return nil
	},
}
