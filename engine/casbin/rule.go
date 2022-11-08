package casbin

import (
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

type PolicyRule struct {
	PType string `json:"p_type,omitempty"`
	V0    string `json:"v0,omitempty"`
	V1    string `json:"v1,omitempty"`
	V2    string `json:"v2,omitempty"`
	V3    string `json:"v3,omitempty"`
	V4    string `json:"v4,omitempty"`
	V5    string `json:"v5,omitempty"`
}

func (line PolicyRule) LoadPolicyLine(model model.Model) error {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}
	return persist.LoadPolicyLine(lineText, model)
}
