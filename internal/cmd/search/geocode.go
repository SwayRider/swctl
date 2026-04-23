package search

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var Geocode = &cli.Command{
	Name:    "geocode",
	Aliases: []string{"gc"},
	Usage:   "Geocode a text query to coordinates",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "query",
			UsageText: "<query> The text to geocode",
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
		&cli.IntFlag{
			Name:  "size",
			Usage: "Maximum number of results",
			Value: 10,
		},
		&cli.StringFlag{
			Name:  "lang",
			Usage: "Preferred language for results (e.g. en, nl)",
		},
		&cli.Float64Flag{
			Name:  "lat",
			Usage: "Focus point latitude",
		},
		&cli.Float64Flag{
			Name:  "lon",
			Usage: "Focus point longitude",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		query := c.StringArg("query")
		if query == "" {
			return fmt.Errorf("query is required")
		}

		results, err := logic.Geocode(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("search-host"),
			c.Int("search-port"),
			c.String("user"),
			c.String("password"),
			query,
			int32(c.Int("size")),
			c.String("lang"),
			c.Float64("lat"),
			c.Float64("lon"),
		)
		if err != nil {
			return err
		}

		if len(results) == 0 {
			fmt.Println("No results found.")
			return nil
		}
		for _, r := range results {
			r.Display()
		}
		return nil
	},
}
