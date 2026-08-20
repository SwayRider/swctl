package prompt

import (
	"context"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
)

func withFakeTerminal(t *testing.T, terminal bool, readFn func(fd int) ([]byte, error)) {
	t.Helper()

	origIsTerminal, origReadPassword := isTerminal, readPassword
	isTerminal = func(fd int) bool { return terminal }
	readPassword = readFn
	t.Cleanup(func() {
		isTerminal = origIsTerminal
		readPassword = origReadPassword
	})
}

func TestPassword_NonTTY_ReturnsErrorWithoutReading(t *testing.T) {
	readCalls := 0
	withFakeTerminal(t, false, func(fd int) ([]byte, error) {
		readCalls++
		return nil, nil
	})

	_, err := Password("Password")
	if err == nil {
		t.Fatal("expected an error on a non-interactive terminal")
	}
	if !strings.Contains(err.Error(), "not an interactive terminal") {
		t.Errorf("expected error to mention non-interactive terminal, got: %v", err)
	}
	if readCalls != 0 {
		t.Errorf("expected ReadPassword not to be called, got %d calls", readCalls)
	}
}

func TestPassword_TTY_ReturnsTypedValue(t *testing.T) {
	withFakeTerminal(t, true, func(fd int) ([]byte, error) {
		return []byte("s3cret"), nil
	})

	got, err := Password("Password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "s3cret" {
		t.Errorf("expected %q, got %q", "s3cret", got)
	}
}

func TestPassword_TTY_EmptyInputIsRequiredError(t *testing.T) {
	withFakeTerminal(t, true, func(fd int) ([]byte, error) {
		return []byte{}, nil
	})

	_, err := Password("Password")
	if err == nil {
		t.Fatal("expected an error for empty input")
	}
	if !strings.Contains(err.Error(), "is required") {
		t.Errorf("expected error to mention it's required, got: %v", err)
	}
}

func TestBeforeFillPassword_AlreadySet_DoesNotPrompt(t *testing.T) {
	readCalls := 0
	withFakeTerminal(t, true, func(fd int) ([]byte, error) {
		readCalls++
		return []byte("unused"), nil
	})

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "password"},
			&cli.StringFlag{Name: "user"},
		},
	}
	if err := cmd.Set("password", "already-set"); err != nil {
		t.Fatalf("failed to seed password flag: %v", err)
	}

	if _, err := BeforeFillPassword(context.Background(), cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if readCalls != 0 {
		t.Errorf("expected no prompt when password already set, got %d read calls", readCalls)
	}
}

func TestBeforeFillPassword_Empty_NonTTY_ReturnsError(t *testing.T) {
	withFakeTerminal(t, false, func(fd int) ([]byte, error) {
		return nil, nil
	})

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "password"},
			&cli.StringFlag{Name: "user"},
		},
	}

	if _, err := BeforeFillPassword(context.Background(), cmd); err == nil {
		t.Fatal("expected an error when password is empty and stdin is not a terminal")
	}
}

func TestBeforeFillPassword_Empty_TTY_FillsPasswordFlag(t *testing.T) {
	withFakeTerminal(t, true, func(fd int) ([]byte, error) {
		return []byte("typed-pass"), nil
	})

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "password"},
			&cli.StringFlag{Name: "user"},
		},
	}

	if _, err := BeforeFillPassword(context.Background(), cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := cmd.String("password"); got != "typed-pass" {
		t.Errorf("expected password flag to be filled with %q, got %q", "typed-pass", got)
	}
}
