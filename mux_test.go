package clir_test

import (
	"testing"

	"maragu.dev/is"

	"maragu.dev/clir"
)

func TestRoute(t *testing.T) {
	t.Run("can handle a root command", func(t *testing.T) {
		mux := clir.NewCommandMux()

		var called bool
		mux.HandleFunc("", func(ctx clir.Context) error {
			called = true
			return nil
		})

		err := mux.Run(clir.Context{
			Args: []string{},
		})
		is.NotError(t, err)
		is.True(t, called)
	})

	t.Run("errors if there is no root command", func(t *testing.T) {
		mux := clir.NewCommandMux()

		err := mux.Run(clir.Context{
			Args: []string{},
		})
		is.Error(t, err, clir.ErrorNotFound)
	})

	t.Run("can handle a subcommand", func(t *testing.T) {
		mux := clir.NewCommandMux()

		var called bool
		mux.HandleFunc("party", func(ctx clir.Context) error {
			called = true
			is.Equal(t, 0, len(ctx.Args))
			return nil
		})

		err := mux.Run(clir.Context{
			Args: []string{"party"},
		})
		is.NotError(t, err)
		is.True(t, called)
	})

	t.Run("errors if there is no subcommand", func(t *testing.T) {
		mux := clir.NewCommandMux()

		err := mux.Run(clir.Context{
			Args: []string{"party"},
		})
		is.Error(t, err, clir.ErrorNotFound)
	})
}
