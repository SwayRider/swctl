package flags

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/urfave/cli/v3"
)

func Optional[T cli.Flag](f T) T {
	if any(f) == nil {
		return f
	}
	return setRequiredProp(f, false)
}

func Required[T cli.Flag](f T) T {
	if any(f) == nil {
		return f
	}
	return setRequiredProp(f, true)
}

func setRequiredProp[T cli.Flag](f T, value bool) T {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return f
	}
	structValue := v.Elem()
	field := structValue.FieldByName("Required")

	if field.IsValid() && field.CanSet() && field.Kind() == reflect.Bool {
		field.SetBool(value)
	}
	return f
}

func Host(serviceName string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:  fmt.Sprintf("%s-host", strings.ToLower(serviceName)),
		Usage: fmt.Sprintf("The host of the %s service", serviceName),
		Sources: cli.EnvVars(
			fmt.Sprintf("%s_HOST", strings.ToUpper(serviceName))),
	}
}

func Port(serviceName string) *cli.IntFlag {
	return &cli.IntFlag{
		Name:  fmt.Sprintf("%s-port", strings.ToLower(serviceName)),
		Usage: fmt.Sprintf("The port of the %s service", serviceName),
		Sources: cli.EnvVars(
			fmt.Sprintf("%s_PORT", strings.ToUpper(serviceName))),
	}
}

func User(env ...string) *cli.StringFlag {
	flag := &cli.StringFlag{
		Name:    "user",
		Aliases: []string{"u"},
		Usage:   "The user to use for authentication",
	}
	if len(env) > 0 {
		flag.Sources = cli.EnvVars(env...)
	}
	return flag
}

func Password(env ...string) *cli.StringFlag {
	flag := &cli.StringFlag{
		Name:    "password",
		Aliases: []string{"p"},
		Usage:   "The password to use for authentication",
	}
	if len(env) > 0 {
		flag.Sources = cli.EnvVars(env...)
	}
	return flag
}
