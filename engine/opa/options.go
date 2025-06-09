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

func WithModulesFromFiles(modules map[string]string) OptFunc {
	return func(s *State) {
		_ = s.InitModulesFromFiles(modules)
	}
}

func WithModulesFromString(modules map[string]string) OptFunc {
	return func(s *State) {
		_ = s.InitModulesFromString(modules)
	}
}

func WithProjectsAuthorizedQuery(query string) OptFunc {
	return func(s *State) {
		_ = s.ParseProjectsQuery(query)
	}
}

func WithFilterAuthorizedPairsQuery(query string) OptFunc {
	return func(s *State) {
		_ = s.ParseFilterPairsQuery(query)
	}
}

func WithFilterAuthorizedProjectsQuery(query string) OptFunc {
	return func(s *State) {
		_ = s.ParseFilterProjectsQuery(query)
	}
}
