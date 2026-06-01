package auth

import (
	"context"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var ListInvites = &cli.Command{
	Name:    "list-invites",
	Aliases: []string{"li"},
	Usage:   "List registration invites",
	Flags: []cli.Flag{
		flags.Required(flags.User("AUTH_USER")),
		flags.Required(flags.Password("AUTH_PASSWORD")),
		&cli.IntFlag{
			Name:  "page",
			Usage: "Page number (0 = all)",
			Value: 0,
		},
		&cli.IntFlag{
			Name:  "page-size",
			Usage: "Number of results per page (0 = all)",
			Value: 0,
		},
		&cli.BoolFlag{
			Name:  "registered",
			Usage: "Filter by registration status: true = registered only, false = pending only (omit for all)",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		var registered *bool
		if c.IsSet("registered") {
			v := c.Bool("registered")
			registered = &v
		}

		invites, err := logic.ListInvites(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
			c.Int("page"),
			c.Int("page-size"),
			registered,
		)
		if err != nil {
			return err
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		t.AppendHeader(table.Row{"Email", "Registered", "ID", "Created At"})

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, WidthMax: 40, WidthMaxEnforcer: text.WrapSoft},
			{Number: 2, WidthMax: 10},
			{Number: 3, WidthMax: 36},
			{Number: 4, WidthMax: 25},
		})

		for _, inv := range invites {
			t.AppendRow(table.Row{
				inv.Email(),
				inv.Registered(),
				inv.Id(),
				inv.CreatedAt().Format(time.RFC3339),
			})
		}
		t.Render()

		return nil
	},
}
