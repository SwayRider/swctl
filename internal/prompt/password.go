// Package prompt provides masked interactive password entry for swctl
// commands, as a fallback when a password wasn't supplied via a positional
// argument, flag, or environment variable.
package prompt

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

// overridable in tests
var (
	isTerminal   = term.IsTerminal
	readPassword = term.ReadPassword
)

// Password prompts on the real terminal for masked input, labeled with
// label. It fails immediately, without prompting, if stdin is not an
// interactive terminal (piped input, CI, redirected from a file), rather
// than risk hanging.
func Password(label string) (string, error) {
	fd := int(os.Stdin.Fd())
	if !isTerminal(fd) {
		return "", fmt.Errorf(
			"%s not provided and stdin is not an interactive terminal; "+
				"pass it as an argument, flag, or environment variable", label)
	}

	fmt.Fprintf(os.Stderr, "%s: ", label)
	b, err := readPassword(fd)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", label, err)
	}
	if len(b) == 0 {
		return "", fmt.Errorf("%s is required", label)
	}
	return string(b), nil
}

// BeforeFillPassword is a cli.BeforeFunc that fills in the "password" flag
// by prompting for masked input when it wasn't supplied via flag or
// environment variable. It runs before urfave/cli's required-flag check, so
// setting the flag here satisfies flags.Required(flags.Password(...)).
func BeforeFillPassword(ctx context.Context, c *cli.Command) (context.Context, error) {
	if c.String("password") != "" {
		return ctx, nil
	}

	label := "Password"
	if u := c.String("user"); u != "" {
		label = "Password for " + u
	}

	pwd, err := Password(label)
	if err != nil {
		return ctx, err
	}
	return ctx, c.Set("password", pwd)
}
