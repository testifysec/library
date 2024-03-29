//go:generate protoc --go_out=. --go_opt=paths=source_relative example.proto

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/testifysec/library/base"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto" // Add this import to handle Protobuf messages
	"google.golang.org/protobuf/reflect/protoreflect"
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
			{
				Name:    "crd",
				Aliases: []string{"c"},
				Usage:   "Print the CRD",
				Action: func(c *cli.Context) error {
					crd()
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

func crd() {
	// Create an instance of ExampleAttestation message.
	exampleAttestation := &ExampleAttestation{}

	// Generate CRD for the envelope
	crd, err := GenerateCRDFromProto(exampleAttestation)
	if err != nil {
		fmt.Println("Error generating CRD:", err)
		return
	}

	fmt.Println(crd)
}

func schema() {
	// Create an instance of ExampleAttestation message.
	exampleAttestation := &ExampleAttestation{}

	// Wrap the ExampleAttestation message in an Any type.

	// Generate JSON schema for the envelope
	jsonSchema, err := generateJsonSchemaFromProto(exampleAttestation)
	if err != nil {
		fmt.Println("Error generating JSON schema:", err)
		return
	}

	fmt.Println(jsonSchema)
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

func generateJsonSchemaFromProto(msg proto.Message) (string, error) {
	return generateJsonSchemaFromProtoReflect(msg.ProtoReflect().Descriptor(), make(map[string]interface{}))
}

func generateJsonSchemaFromProtoReflect(desc protoreflect.MessageDescriptor, schemas map[string]interface{}) (string, error) {
	title := string(desc.FullName())
	if schema, ok := schemas[title]; ok {
		return schema.(string), nil
	}

	properties := map[string]interface{}{}
	for i := 0; i < desc.Fields().Len(); i++ {
		fd := desc.Fields().Get(i)
		fieldName := string(fd.Name())
		fieldSchema := map[string]interface{}{}

		switch fd.Kind() {
		case protoreflect.BoolKind:
			fieldSchema["type"] = "boolean"
		case protoreflect.EnumKind, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind:
			fieldSchema["type"] = "integer"
		case protoreflect.Fixed32Kind, protoreflect.Fixed64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
			fieldSchema["type"] = "number"
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			fieldSchema["type"] = "number"
		case protoreflect.StringKind:
			fieldSchema["type"] = "string"
		case protoreflect.BytesKind:
			fieldSchema["type"] = "string"
		case protoreflect.MessageKind, protoreflect.GroupKind:
			if fd.Message().FullName() == "google.protobuf.Timestamp" {
				fieldSchema["type"] = "string"
				fieldSchema["format"] = "date-time"
			} else {
				nestedSchema, err := generateJsonSchemaFromProtoReflect(fd.Message(), schemas)
				if err != nil {
					return "", err
				}
				var nestedSchemaObj map[string]interface{}
				if err := json.Unmarshal([]byte(nestedSchema), &nestedSchemaObj); err != nil {
					return "", err
				}
				fieldSchema = nestedSchemaObj
			}
		}

		if fd.IsList() {
			properties[fieldName] = map[string]interface{}{
				"type":  "array",
				"items": fieldSchema,
			}
		} else {
			properties[fieldName] = fieldSchema
		}
	}

	schema := map[string]interface{}{
		"$schema":    "http://json-schema.org/draft-07/schema#",
		"title":      title,
		"type":       "object",
		"properties": properties,
	}

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}

	schemaStr := string(schemaJSON)
	schemas[title] = schemaStr
	return schemaStr, nil
}

// GenerateCRDFromProto creates a CRD directly from a Protobuf message, including the schema.
func GenerateCRDFromProto(msg proto.Message) (string, error) {
	properties := make(map[string]interface{})
	fields := msg.ProtoReflect().Descriptor().Fields()

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		fieldSchema := map[string]interface{}{}
		switch fd.Kind() {
		case protoreflect.BoolKind:
			fieldSchema["type"] = "boolean"
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
			protoreflect.Fixed32Kind, protoreflect.Fixed64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
			fieldSchema["type"] = "integer"
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			fieldSchema["type"] = "number"
		case protoreflect.StringKind:
			fieldSchema["type"] = "string"
		case protoreflect.BytesKind:
			fieldSchema["type"] = "string" // Consider specifying format as "byte" for base64 encoded data
		case protoreflect.EnumKind:
			// Assuming enum values should be represented as strings
			fieldSchema["type"] = "string"
			// Optionally, add "enum" field to list possible values
		case protoreflect.MessageKind, protoreflect.GroupKind:
			// This simplified example assumes sub-messages/groups are not nested.
			fieldSchema["type"] = "object"
		}
		properties[string(fd.Name())] = fieldSchema
	}

	crd := map[string]interface{}{
		"version": "https://witness.dev/attestations/custom/v0.1",
		"metadata": map[string]string{
			"name": "organization",
			"type": "acme.corp.io/compliance",
		},
		"spec": map[string]interface{}{
			"attest": map[string]interface{}{
				"type": "postproduct",
				"executor": map[string]interface{}{
					"type": "command",
					"arguments": []string{
						"./custom-attestor-test/custom-attestor",
						"attest",
					},
				},
			},
			"verify": map[string]interface{}{
				"type": "postproduct",
				"executor": map[string]interface{}{
					"type": "command",
					"arguments": []string{
						"./custom-attestor-test/custom-attestor",
						"verify",
					},
				},
			},
			"versions": []interface{}{
				map[string]interface{}{
					"name": "v0.1",
					"schema": map[string]interface{}{
						"type":       "object",
						"properties": properties,
						// Add required fields, enums, or other constraints as needed
					},
				},
			},
		},
	}

	crdJSON, err := json.MarshalIndent(crd, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal CRD to JSON: %w", err)
	}

	return string(crdJSON), nil
}
