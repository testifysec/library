package sdk

import (
	"github.com/testifysec/library/sdk/plugin"

	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"
)

// RegisterSchemaCommand registers a CLI command for generating a JSON schema.
func registerSchemaCommand(app *cli.App, protoMsg proto.Message) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "schema",
		Usage: "Generate the schema",
		Action: func(c *cli.Context) error {
			schema, err := GenerateJSONSchema(protoMsg)
			if err != nil {
				return err
			}
			c.App.Writer.Write([]byte(schema))
			return nil
		},
	})
}

// RegisterCRDCommand registers a CLI command for generating a CRD.
func registerCRDCommand(app *cli.App, protoMsg proto.Message) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "crd",
		Usage: "Generate the CRD",
		Action: func(c *cli.Context) error {
			crd, err := GenerateCRD(protoMsg)
			if err != nil {
				return err
			}
			c.App.Writer.Write([]byte(crd))
			return nil
		},
	})
}

func registerAttestCommand(app *cli.App, p plugin.Plugin) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:   "attest",
		Usage:  "Execute the attest action",
		Action: func(c *cli.Context) error { return p.Attest() },
	})
}

func registerVerifyCommand(app *cli.App, p plugin.Plugin) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:   "verify",
		Usage:  "Execute the verify action",
		Action: func(c *cli.Context) error { return p.Verify() },
	})
}

func GenerateJSONSchema(protoMsg proto.Message) (string, error) {
	// Implement your logic here
	return "", nil
}

func GenerateCRD(protoMsg proto.Message) (string, error) {
	// Implement your logic here
	return "", nil
}

// single register function
func RegisterCommands(app *cli.App, protoMsg proto.Message, p plugin.Plugin) {
	registerSchemaCommand(app, protoMsg)
	registerCRDCommand(app, protoMsg)
	registerAttestCommand(app, p)
	registerVerifyCommand(app, p)
}
