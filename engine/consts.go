package engine

type Type string

const (
	Noop     Type = "noop"
	Casbin   Type = "casbin"
	Opa      Type = "opa"
	Zanzibar Type = "zanzibar"
)
