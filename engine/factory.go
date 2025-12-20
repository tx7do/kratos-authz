package engine

import (
	"context"
	"fmt"
	"sync"
)

type FactoryFunc func(ctx context.Context, options ...any) (Engine, error)

var (
	factoryMu sync.RWMutex
	factories = make(map[Type]FactoryFunc)
)

// NewEngine creates a new authz Engine based on the registered FactoryFunc for the given Type.
func NewEngine(ctx context.Context, typ Type, options ...any) (Engine, error) {
	factory, ok := GetFactory(typ)
	if !ok {
		return nil, fmt.Errorf("authz: unknown engine type %s", typ)
	}
	return factory(ctx, options...)
}

// Register an authz engine factory function
func Register(typ Type, f FactoryFunc) error {
	factoryMu.Lock()
	defer factoryMu.Unlock()
	if _, ok := factories[typ]; ok {
		return fmt.Errorf("authz: engine factory %s already registered", typ)
	}
	factories[typ] = f
	return nil
}

// GetFactory returns a registered FactoryFunc for a given Type and whether it existed.
// Safe for concurrent use.
func GetFactory(typ Type) (FactoryFunc, bool) {
	factoryMu.RLock()
	defer factoryMu.RUnlock()
	f, ok := factories[typ]
	return f, ok
}

// ListFactories returns a slice of currently registered Types.
func ListFactories() []Type {
	factoryMu.RLock()
	defer factoryMu.RUnlock()
	res := make([]Type, 0, len(factories))
	for k := range factories {
		res = append(res, k)
	}
	return res
}

// Unregister removes a registered factory by Type. It returns true if a factory was removed.
// Use with caution in concurrent environments (primarily intended for tests).
func Unregister(typ Type) bool {
	factoryMu.Lock()
	defer factoryMu.Unlock()
	if _, ok := factories[typ]; ok {
		delete(factories, typ)
		return true
	}
	return false
}
