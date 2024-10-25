package clir_test

import (
	"flag"
	"strings"
	"testing"

	"maragu.dev/is"

	"maragu.dev/clir"
)

func TestRouter_Run(t *testing.T) {
	t.Run("errors on run if there is no root route", func(t *testing.T) {
		r := clir.NewRouter()

		err := r.Run(clir.Context{})
		is.Error(t, err, clir.ErrorRouteNotFound)
	})

	t.Run("errors on run if there is no named route", func(t *testing.T) {
		r := clir.NewRouter()

		err := r.Run(clir.Context{
			Args: []string{"dance"},
		})
		is.Error(t, err, clir.ErrorRouteNotFound)
	})
}

func TestRouter_RouteFunc(t *testing.T) {
	t.Run("can route and run a root route", func(t *testing.T) {
		r := clir.NewRouter()

		var called bool
		r.RouteFunc("", func(ctx clir.Context) error {
			called = true
			return nil
		})

		err := r.Run(clir.Context{})
		is.NotError(t, err)
		is.True(t, called)
	})

	t.Run("can route and run a named route", func(t *testing.T) {
		r := clir.NewRouter()

		var called bool
		r.RouteFunc("dance", func(ctx clir.Context) error {
			called = true
			is.Equal(t, 0, len(ctx.Args))
			return nil
		})

		err := r.Run(clir.Context{
			Args: []string{"dance"},
		})
		is.NotError(t, err)
		is.True(t, called)
	})

	t.Run("panics if the route already exists", func(t *testing.T) {
		r := clir.NewRouter()

		r.RouteFunc("", func(ctx clir.Context) error {
			return nil
		})

		defer func() {
			if rec := recover(); rec == nil {
				t.FailNow()
			}
		}()

		r.RouteFunc("", func(ctx clir.Context) error {
			return nil
		})
	})

	t.Run("supports regular expression in routes", func(t *testing.T) {
		r := clir.NewRouter()

		var called bool
		r.RouteFunc(`\w+`, func(ctx clir.Context) error {
			called = true
			is.Equal(t, 1, len(ctx.Matches))
			is.Equal(t, "dance", ctx.Matches[0])
			return nil
		})

		err := r.Run(clir.Context{
			Args: []string{"dance"},
		})
		is.NotError(t, err)
		is.True(t, called)
	})

	t.Run("supports regular expression in routes including submatches", func(t *testing.T) {
		r := clir.NewRouter()

		var called bool
		r.RouteFunc(`(\w+)`, func(ctx clir.Context) error {
			called = true
			is.Equal(t, 2, len(ctx.Matches))
			is.Equal(t, "dance", ctx.Matches[0])
			is.Equal(t, "dance", ctx.Matches[1])
			return nil
		})

		err := r.Run(clir.Context{
			Args: []string{"dance"},
		})
		is.NotError(t, err)
		is.True(t, called)
	})
}

func TestRouter_Use(t *testing.T) {
	t.Run("can use middlewares on root and named routes", func(t *testing.T) {
		r := clir.NewRouter()

		m1 := newMiddleware(t, "m1")
		m2 := newMiddleware(t, "m2")

		r.Use(m1, m2)

		r.RouteFunc("", func(ctx clir.Context) error {
			ctx.Println("root")
			return nil
		})

		r.RouteFunc("dance", func(ctx clir.Context) error {
			ctx.Println("dance")
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
			Args: []string{"dance"},
			Out:  &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nm2\ndance\n", b.String())
	})

	t.Run("panics if routes are already registered", func(t *testing.T) {
		r := clir.NewRouter()

		r.RouteFunc("", func(ctx clir.Context) error {
			return nil
		})

		defer func() {
			if rec := recover(); rec == nil {
				t.FailNow()
			}
		}()

		r.Use(newMiddleware(t, "m1"))
	})

	t.Run("can use middleware that parses flags", func(t *testing.T) {
		r := clir.NewRouter()

		r.Use(func(next clir.Runner) clir.Runner {
			return clir.RunnerFunc(func(ctx clir.Context) error {
				fs := flag.NewFlagSet("test", flag.ContinueOnError)
				v := fs.Bool("v", false, "")
				err := fs.Parse(ctx.Args)
				is.NotError(t, err)
				is.True(t, *v)

				t.Log(fs.Args())
				ctx.Args = fs.Args()

				return next.Run(ctx)
			})
		})

		var called bool
		r.RouteFunc("", func(ctx clir.Context) error {
			called = true
			return nil
		})

		err := r.Run(clir.Context{
			Args: []string{"-v"},
		})
		is.NotError(t, err)
		is.True(t, called)
	})

	t.Run("does not call route if middleware doesn't call next", func(t *testing.T) {
		r := clir.NewRouter()

		r.Use(func(next clir.Runner) clir.Runner {
			return clir.RunnerFunc(func(ctx clir.Context) error {
				return nil
			})
		})

		var called bool
		r.RouteFunc("", func(ctx clir.Context) error {
			called = true
			return nil
		})

		err := r.Run(clir.Context{})
		is.NotError(t, err)
		is.True(t, !called)
	})
}

//nolint:staticcheck
func TestRouter_Scope(t *testing.T) {
	t.Skip("not implemented")

	t.Run("can scope routes with a new middleware stack", func(t *testing.T) {
		r := clir.NewRouter()

		m1 := newMiddleware(t, "m1")
		m2 := newMiddleware(t, "m2")
		m3 := newMiddleware(t, "m3")

		// Only apply the first one here
		r.Use(m1)

		r.RouteFunc("", func(ctx clir.Context) error {
			ctx.Println("root")
			return nil
		})

		r.Scope(func(r *clir.Router) {
			r.Use(m2)

			r.RouteFunc("dance", func(ctx clir.Context) error {
				ctx.Println("dance")
				return nil
			})
		})

		r.Scope(func(r *clir.Router) {
			r.Use(m3)

			r.RouteFunc("sleep", func(ctx clir.Context) error {
				ctx.Println("sleep")
				return nil
			})
		})

		var b strings.Builder
		err := r.Run(clir.Context{
			Out: &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nroot\n", b.String())

		b.Reset()

		err = r.Run(clir.Context{
			Args: []string{"dance"},
			Out:  &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nm2\ndance\n", b.String())

		b.Reset()

		err = r.Run(clir.Context{
			Args: []string{"sleep"},
			Out:  &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nm3\nsleep\n", b.String())
	})

	t.Run("can nest scopes", func(t *testing.T) {
		r := clir.NewRouter()

		m1 := newMiddleware(t, "m1")
		m2 := newMiddleware(t, "m2")
		m3 := newMiddleware(t, "m3")

		r.Use(m1)

		r.RouteFunc("", func(ctx clir.Context) error {
			ctx.Println("root")
			return nil
		})

		r.Scope(func(r *clir.Router) {
			r.Use(m2)

			r.RouteFunc("dance", func(ctx clir.Context) error {
				ctx.Println("dance")
				return nil
			})

			r.Scope(func(r *clir.Router) {
				r.Use(m3)

				r.RouteFunc("sleep", func(ctx clir.Context) error {
					ctx.Println("sleep")
					return nil
				})
			})
		})

		var b strings.Builder
		err := r.Run(clir.Context{
			Out: &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nroot\n", b.String())

		b.Reset()

		err = r.Run(clir.Context{
			Args: []string{"dance"},
			Out:  &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nm2\ndance\n", b.String())

		b.Reset()

		err = r.Run(clir.Context{
			Args: []string{"sleep"},
			Out:  &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nm2\nm3\nsleep\n", b.String())
	})
}

func TestRouter_Branch(t *testing.T) {
	t.Run("can branch into a new router with a new middleware stack", func(t *testing.T) {
		r := clir.NewRouter()

		m1 := newMiddleware(t, "m1")
		m2 := newMiddleware(t, "m2")

		r.Use(m1)

		r.Branch("dance", func(r *clir.Router) {
			r.Use(m2)

			r.RouteFunc("", func(ctx clir.Context) error {
				ctx.Println("dance root")
				return nil
			})

			r.RouteFunc("list", func(ctx clir.Context) error {
				ctx.Println("dance list")
				return nil
			})
		})

		var b strings.Builder
		err := r.Run(clir.Context{
			Args: []string{"dance"},
			Out:  &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nm2\ndance root\n", b.String())

		b.Reset()

		err = r.Run(clir.Context{
			Args: []string{"dance", "list"},
			Out:  &b,
		})
		is.NotError(t, err)
		is.Equal(t, "m1\nm2\ndance list\n", b.String())
	})

	t.Run("panics if the route already exists", func(t *testing.T) {
		r := clir.NewRouter()

		r.RouteFunc("", func(ctx clir.Context) error {
			return nil
		})

		defer func() {
			if rec := recover(); rec == nil {
				t.FailNow()
			}
		}()

		r.Branch("", func(r *clir.Router) {})
	})
}

func newMiddleware(t *testing.T, name string) clir.Middleware {
	t.Helper()

	return func(next clir.Runner) clir.Runner {
		return clir.RunnerFunc(func(ctx clir.Context) error {
			ctx.Println(name)
			return next.Run(ctx)
		})
	}
}
