package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var EnsureServiceClient = &cli.Command{
	Name:    "ensure-service-client",
	Aliases: []string{"esc"},
	Usage:   "Idempotently create a service client and write credentials to a file",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "name",
			UsageText: "<name>",
		},
		&cli.StringArgs{
			Name:      "scope",
			UsageText: "<scope...>",
			Min:       1,
			Max:       -1,
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
		&cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o"},
			Usage:    "Write credentials to this file in KEY=VALUE format",
			Required: true,
		},
		&cli.IntFlag{
			Name:  "retries",
			Usage: "Number of retries on connection failure (with 3s delay between each)",
			Value: 10,
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		output := c.String("output")

		// Idempotency: skip if credentials file already exists and is non-empty.
		if fi, err := os.Stat(output); err == nil && fi.Size() > 0 {
			fmt.Printf("credentials already present at %s, skipping registration\n", output)
			return nil
		}

		name := c.StringArg("name")
		retries := c.Int("retries")

		var clientId, clientSecret string
		var err error
		for i := range retries + 1 {
			clientId, clientSecret, err = logic.CreateServiceClient(
				c.String("auth-host"),
				c.Int("auth-port"),
				c.String("user"),
				c.String("password"),
				name,
				"",
				c.StringArgs("scope"),
			)
			if err == nil {
				break
			}
			if isAlreadyExists(err) {
				return fmt.Errorf(
					"service client %q already exists in authservice but no credentials file found at %s.\n"+
						"Recovery: run `swctl auth list-service-clients` to find the client ID, "+
						"then `swctl auth delete-service-client <id>` and retry",
					name, output)
			}
			if i < retries {
				fmt.Printf("authservice not ready (attempt %d/%d), retrying in 3s...\n", i+1, retries+1)
				time.Sleep(3 * time.Second)
			}
		}
		if err != nil {
			return fmt.Errorf("failed to create service client after %d attempts: %w", retries+1, err)
		}

		content := fmt.Sprintf("SWAYRIDER_API_CLIENT_ID=%s\nSWAYRIDER_API_CLIENT_SECRET=%s\n", clientId, clientSecret)
		if err := os.WriteFile(output, []byte(content), 0600); err != nil {
			return fmt.Errorf("failed to write credentials to %s: %w", output, err)
		}

		fmt.Printf("service client %q registered; credentials written to %s\n", name, output)
		return nil
	},
}

func isAlreadyExists(err error) bool {
	s, ok := status.FromError(err)
	return ok && s.Code() == codes.AlreadyExists
}
