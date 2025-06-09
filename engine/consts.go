package engine

type Type string

const (
	Casbin   Type = "casbin"
	Opa      Type = "opa"
	Zanzibar Type = "zanzibar"
)
