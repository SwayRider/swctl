package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/logic"
)

var CheckPasswordStrength = &cli.Command{
	Name:    "check-password-strength",
	Aliases: []string{"cps"},
	Usage:   "Check the strength of a password",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "password",
			UsageText: "<password> The password to check",
		},
	},
	Flags: []cli.Flag{
		//flags.Required(flags.Password()),
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		pwd := c.StringArg("password")
		if pwd == "" {
			return fmt.Errorf("password is required")
		}

		_, message, err := logic.CheckPasswordStrength(
			c.String("auth-host"),
			c.Int("auth-port"),
			pwd,
		)
		if err != nil {
			return err
		}

		fmt.Println(message)
		return nil
	},
}
