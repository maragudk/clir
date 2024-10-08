package clir_test

import (
	"strings"
	"testing"

	"maragu.dev/is"

	"maragu.dev/clir"
)

func TestCommandRouter_Run(t *testing.T) {
	t.Run("can route and run a root command", func(t *testing.T) {
		r := clir.NewCommandRouter()

		var called bool
		r.SubFunc("", func(ctx clir.Context) error {
			called = true
			return nil
		})

		err := r.Run(clir.Context{})
		is.NotError(t, err)
		is.True(t, called)
	})

	t.Run("errors on run if there is no root command", func(t *testing.T) {
		r := clir.NewCommandRouter()

		err := r.Run(clir.Context{})
		is.Error(t, err, clir.ErrorNotFound)
	})

	t.Run("can route and run a subcommand", func(t *testing.T) {
		r := clir.NewCommandRouter()

		var called bool
		r.SubFunc("party", func(ctx clir.Context) error {
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

	t.Run("can use middlewares on root and subcommands", func(t *testing.T) {
		r := clir.NewCommandRouter()

		m1 := newMiddleware(t, "m1")
		m2 := newMiddleware(t, "m2")

		r.Use(m1, m2)

		r.SubFunc("", func(ctx clir.Context) error {
			ctx.Println("root")
			return nil
		})

		r.SubFunc("party", func(ctx clir.Context) error {
			ctx.Println("party")
			return nil
		})

		var b strings.Builder
		err := r.Run(clir.Context{
			Out: &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nm2\nroot\n", b.String())

		b.Reset()

		err = r.Run(clir.Context{
			Args: []string{"party"},
			Out:  &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nm2\nparty\n", b.String())
	})
}

func newMiddleware(t *testing.T, name string) clir.Middleware {
	return func(next clir.Command) clir.Command {
		return clir.CommandFunc(func(ctx clir.Context) error {
			ctx.Println(name)
			return next.Run(ctx)
		})
	}
}
