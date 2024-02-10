// plugin/interface.go
package plugin

type Plugin interface {
	Attest() error
	Verify() error
}
