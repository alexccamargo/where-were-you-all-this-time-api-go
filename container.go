//+build wireinject

package main

import (
	"github.com/alexccamargo/wwuatt/services"
	"github.com/google/wire"
)

// InitializeDI initialize dependencies
func InitializeDI() *services.MovieService {
	panic(wire.Build(
		services.NewConnection,
		services.NewMovieService,
	))
	return &services.MovieService{}
}
