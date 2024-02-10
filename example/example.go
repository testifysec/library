//go:generate protoc --go_out=. --go_opt=paths=source_relative example.proto

package main

import (
	"encoding/json"
	"fmt"
	"os"
	reflect "reflect"

	"github.com/testifysec/library/base"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"             // Add this import to handle Protobuf messages
	"google.golang.org/protobuf/types/known/anypb" // For wrapping messages in Any
	"google.golang.org/protobuf/types/known/structpb"
)

func main() {
	app := &cli.App{
		Name:  "mycli",
		Usage: "A simple CLI with subcommands",
		Commands: []*cli.Command{
			{
				Name:    "attest",
				Aliases: []string{"a"},
				Usage:   "Attest a statement",
				Action: func(c *cli.Context) error {
					attest()
					return nil
				},
			},
			{
				Name:    "verify",
				Aliases: []string{"v"},
				Usage:   "Verify a statement",
				Action: func(c *cli.Context) error {
					verify()
					return nil
				},
			},
			{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Print the schema",
				Action: func(c *cli.Context) error {
					schema()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func schema() {

	generateSchema(&ExampleAttestation{})

}

func attest() {
	// Use the generated OrganizationAttestation message
	attestation := &ExampleAttestation{
		Name: "MyOrg",
		Team: "MyTeam",
	}

	// Wrap the attestation in the AttestationWorkflowEnvelope
	envelope := &base.AttestationWorkflowEnvelope{
		Metadata: &base.Metadata{},
		Verify:   &base.Executor{},
		Attest:   &base.Executor{Type: base.ExecutorType_COMMAND, Arguments: []string{"echo", "hello"}},
		Schema:   &structpb.Struct{},
		Payload:  mustWrapAny(attestation),
	}

	// Marshal the envelope to JSON
	var jsonData []byte
	var err error
	jsonData, err = json.Marshal(envelope)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		os.Exit(1)
	}

	// Print the JSON data to standard output
	fmt.Print(string(jsonData))
}

func verify() {
	fmt.Println("unimplemented")
}

// mustWrapAny is a helper function to wrap a message in a google.protobuf.Any
// It panics if the operation fails, which is suitable for static data like this example.
func mustWrapAny(msg proto.Message) *anypb.Any {
	any, err := anypb.New(msg)
	if err != nil {
		panic(fmt.Sprintf("Failed to wrap message in Any: %v", err))
	}
	return any
}

func generateSchema(msg proto.Message) (string, error) {
	// This is a very simplified and static approach to generate a JSON Schema for demonstration.
	// A real implementation would need to inspect `msg` dynamically using protobuf reflection.

	schema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   reflect.TypeOf(msg).Elem().Name(),
		"type":    "object",
		"properties": map[string]interface{}{
			// Assuming these are the fields in your attestation protobuf message.
			// Replace these with dynamic inspection of `msg` for a real implementation.
			"name": map[string]string{
				"type": "string",
			},
			"age": map[string]string{
				"type": "integer",
			},
		},
		// Assuming the `name` field is required.
		"required": []string{"name"},
	}

	// You would wrap this schema inside the structure of AttestationWorkflowEnvelope's schema
	// For brevity, this example focuses on generating schema for the message itself.

	jsonSchema, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonSchema), nil
}
