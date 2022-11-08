package opa

import (
	"github.com/open-policy-agent/opa/ast"
)

type OptFunc func(*State)

func WithModules(mods map[string]*ast.Module) OptFunc {
	return func(s *State) {
		s.modules = mods
	}
}
