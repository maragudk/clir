package clir_test

import (
	"testing"

	"maragu.dev/is"

	"maragu.dev/clir"
)

func TestRoute(t *testing.T) {
	t.Run("can route a subcommand", func(t *testing.T) {
		r := clir.NewRouter()

		var called bool
		r.Route("sub", func(ctx clir.Context) error {
			called = true
			return nil
		})

		clir.Main([]string{"clir", "sub"}, nil, nil, nil, r)

		is.True(t, called)
	})
}
