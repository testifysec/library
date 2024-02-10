// plugin/interface.go
package plugin

type Plugin interface {
	Attest() error
	Verify() error
	GenerateSchema() error
	GenerateCRD() error
}
