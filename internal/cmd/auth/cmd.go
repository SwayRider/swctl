package auth

import (
	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
)

var Command = &cli.Command{
	Name:  "auth",
	Usage: "Interact with the authservice",
	Flags: []cli.Flag{
		flags.Required(flags.Host("Auth")),
		flags.Required(flags.Port("Auth")),
	},
	Commands: []*cli.Command{
		CheckPasswordStrength,
		CreateAdmin,
		CreateUser,
		CreateServiceClient,
		ListServiceClients,
		DeleteServiceClient,
	},
}
