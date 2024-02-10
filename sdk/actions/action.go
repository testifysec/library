package actions

type Attester interface {
	Attest() error
}

type Verifier interface {
	Verify() error
}

type SchemaGenerator interface {
	GenerateSchema() error
}

type CRDGenerator interface {
	GenerateCRD() error
}
