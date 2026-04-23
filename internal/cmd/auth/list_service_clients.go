package auth

import (
	"context"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v3"
	"github.com/swayrider/swctl/internal/flags"
	"github.com/swayrider/swctl/internal/logic"
)

var ListServiceClients = &cli.Command{
	Name:    "list-service-clients",
	Aliases: []string{"lsc"},
	Usage:   "List all service clients",
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
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		clients, err := logic.ListServiceClients(
			c.String("auth-host"),
			c.Int("auth-port"),
			c.String("user"),
			c.String("password"),
			c.Int("page"),
			c.Int("page-size"),
		)
		if err != nil {
			return err
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		t.AppendHeader(table.Row{"Name", "ClientID", "Scopes", "Description"})

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, WidthMax: 20, WidthMaxEnforcer: text.WrapSoft},
			{Number: 2, WidthMax: 32},
			{Number: 3, WidthMax: 16, WidthMaxEnforcer: text.WrapSoft},
			{Number: 4, WidthMax: 64, WidthMaxEnforcer: text.WrapSoft},
		})

		for _, cln := range clients {
			scopeText := strings.Join(cln.Scopes(), " ")
			t.AppendRow(table.Row{cln.Name(), cln.ClientId(), scopeText, cln.Description()})
		}
		t.Render()

		return nil
	},
}
