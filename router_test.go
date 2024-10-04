package clir_test

import (
	"testing"

	"maragu.dev/is"

	"maragu.dev/clir"
)

func TestCommandRouter_Run(t *testing.T) {
	t.Run("can route and run a root command", func(t *testing.T) {
		r := clir.NewCommandRouter()

		var called bool
		r.RouteFunc("", func(ctx clir.Context) error {
			called = true
			return nil
		})

		err := r.Run(clir.Context{
			Args: []string{},
		})
		is.NotError(t, err)
		is.True(t, called)
	})

	t.Run("errors on run if there is no root command", func(t *testing.T) {
		r := clir.NewCommandRouter()

		err := r.Run(clir.Context{
			Args: []string{},
		})
		is.Error(t, err, clir.ErrorNotFound)
	})

	t.Run("can route and run a subcommand", func(t *testing.T) {
		r := clir.NewCommandRouter()

		var called bool
		r.RouteFunc("party", func(ctx clir.Context) error {
			called = true
			is.Equal(t, 0, len(ctx.Args))
			return nil
		})

		err := r.Run(clir.Context{
			Args: []string{"party"},
		})
		is.NotError(t, err)
		is.True(t, called)
	})

	t.Run("errors on run if there is no subcommand", func(t *testing.T) {
		r := clir.NewCommandRouter()

		err := r.Run(clir.Context{
			Args: []string{"party"},
		})
		is.Error(t, err, clir.ErrorNotFound)
	})
}
