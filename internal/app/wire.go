//go:build wireinject

package app

import "github.com/google/wire"

////////////////////////////////////////////////////////////////////////////////////
// This file is a stub for Wire injector, use "make wire" to rebuild dependencies //
////////////////////////////////////////////////////////////////////////////////////

// BuildService constructs service instance with app dependencies using Wire.
func BuildService() (*Service, error) {
	wire.Build(dependenciesSet)
	return &Service{}, nil
}
