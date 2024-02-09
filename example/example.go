package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

type Predicate struct {
	Type        string       `json:"type"`
	Attestation Organization `json:"attestation"`
}

type Organization struct {
	Name  string `json:"name"`
	Team  string `json:"team"`
	Cloud string `json:"cloud" enum:"aws,azure,gcp"`
}

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
	attestation := Predicate{
		Type: "acme.corp.io/compliance/organization/v0.1",
		Attestation: Organization{
			Name:  "MyOrg",
			Team:  "MyTeam",
			Cloud: "foobar",
		},
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(attestation)
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
