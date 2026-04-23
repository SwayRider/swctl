package search

import (
	"context"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var Reverse = &cli.Command{
	Name:    "reverse",
	Aliases: []string{"rv"},
	Usage:   "Reverse geocode coordinates to a place name",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "lat",
			UsageText: "<lat> Latitude",
		},
		&cli.StringArg{
			Name:      "lon",
			UsageText: "<lon> Longitude",
		},
	},
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
		&cli.IntFlag{
			Name:  "size",
			Usage: "Maximum number of results",
			Value: 5,
		},
		&cli.StringFlag{
			Name:  "lang",
			Usage: "Preferred language for results (e.g. en, nl)",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		latStr := c.StringArg("lat")
		lonStr := c.StringArg("lon")
		if latStr == "" || lonStr == "" {
			return fmt.Errorf("lat and lon are required")
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return fmt.Errorf("invalid lat: %w", err)
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return fmt.Errorf("invalid lon: %w", err)
		}

		results, err := logic.ReverseGeocode(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("search-host"),
			c.Int("search-port"),
			c.String("user"),
			c.String("password"),
			lat,
			lon,
			int32(c.Int("size")),
			c.String("lang"),
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
