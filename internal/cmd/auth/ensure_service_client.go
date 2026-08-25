package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
	"github.com/swayrider/swctl/internal/prompt"
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
	Before: prompt.BeforeFillPassword,
	Action: func(ctx context.Context, c *cli.Command) error {
		output := c.String("output")
		name := c.StringArg("name")
		retries := c.Int("retries")
		authHost := c.String("auth-host")
		authPort := c.Int("auth-port")
		user := c.String("user")
		password := c.String("password")

		// Idempotency: a local credentials file alone isn't enough to skip --
		// it survives things the file knows nothing about (a reset/fresh
		// database, or being pointed at a different authservice), leaving
		// stale credentials that no longer exist server-side. Only skip once
		// the target authservice confirms THIS SPECIFIC client ID is still
		// there -- matching by name alone isn't enough either, since two
		// different authservice instances can each have their own same-named
		// "swayrider-api" client with different IDs/secrets.
		if fi, err := os.Stat(output); err == nil && fi.Size() > 0 {
			existingId, parseErr := readClientId(output)
			var exists bool
			if parseErr == nil {
				exists, err = serviceClientIdExistsWithRetries(authHost, authPort, user, password, existingId, retries)
				if err != nil {
					return fmt.Errorf("failed to check for existing service client %q after %d attempts: %w", name, retries+1, err)
				}
			}
			if exists {
				fmt.Printf("credentials already present at %s and service client %s still exists in authservice, skipping registration\n", output, existingId)
				return nil
			}
			fmt.Printf("credentials file at %s is stale (client not found in authservice -- reset database, or different authservice?) -- re-registering\n", output)
		}

		var clientId, clientSecret string
		var err error
		for i := range retries + 1 {
			clientId, clientSecret, err = logic.CreateServiceClient(
				authHost,
				authPort,
				user,
				password,
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

// readClientId extracts SWAYRIDER_API_CLIENT_ID from a credentials file
// previously written by this command (KEY=VALUE per line).
func readClientId(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), "=")
		if found && key == "SWAYRIDER_API_CLIENT_ID" && value != "" {
			return value, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("SWAYRIDER_API_CLIENT_ID not found in %s", path)
}

// serviceClientIdExistsWithRetries reports whether a service client with the
// given ID currently exists in the target authservice, retrying (like the
// create path) while authservice is still booting. A large pageSize is used
// instead of paging through results -- a local/dev authservice only ever
// holds a handful of service clients, so this is simpler than it's worth
// making exhaustive.
func serviceClientIdExistsWithRetries(
	authHost string, authPort int, user, password, clientId string, retries int,
) (bool, error) {
	var clients []*logic.ServiceClient
	var err error
	for i := range retries + 1 {
		clients, err = logic.ListServiceClients(authHost, authPort, user, password, 1, 100)
		if err == nil {
			break
		}
		if i < retries {
			fmt.Printf("authservice not ready (attempt %d/%d), retrying in 3s...\n", i+1, retries+1)
			time.Sleep(3 * time.Second)
		}
	}
	if err != nil {
		return false, err
	}
	for _, sc := range clients {
		if sc.ClientId() == clientId {
			return true, nil
		}
	}
	return false, nil
}
