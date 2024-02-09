package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/testifysec/library/base"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"             // Add this import to handle Protobuf messages
	"google.golang.org/protobuf/types/known/anypb" // For wrapping messages in Any
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func attest() {
	// Use the generated OrganizationAttestation message
	attestation := &OrganizationAttestation{
		Name:  "MyOrg",
		Team:  "MyTeam",
		Cloud: CloudProvider_AWS, // Example, replace "AWS" with actual value
	}

	// Wrap the attestation in the AttestationWorkflowEnvelope
	envelope := &base.AttestationWorkflowEnvelope{
		Metadata: &base.Metadata{},
		Verify:   &base.Executor{},
		Attest: &base.Executor{
			Type:      base.ExecutorType_COMMAND,
			Arguments: []string{"echo", "hello"},
		},
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
