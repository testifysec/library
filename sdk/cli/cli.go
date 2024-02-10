package cli

import (
	"github.com/testifysec/library/sdk/plugin"

	favecli "github.com/urfave/cli/v2"
)

func RegisterCommands(app *favecli.App, p plugin.Plugin) {
	app.Commands = []*favecli.Command{
		{
			Name:   "attest",
			Usage:  "Execute the attest action",
			Action: func(c *favecli.Context) error { return p.Attest() },
		},
		{
			Name:   "verify",
			Usage:  "Execute the verify action",
			Action: func(c *favecli.Context) error { return p.Verify() },
		},
		{
			Name:   "schema",
			Usage:  "Generate the schema",
			Action: func(c *favecli.Context) error { return p.GenerateSchema() },
		},
		{
			Name:   "crd",
			Usage:  "Generate the CRD",
			Action: func(c *favecli.Context) error { return p.GenerateCRD() },
		},
	}
}
