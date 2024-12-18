package middleware_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"maragu.dev/is"

	"maragu.dev/clir"
	"maragu.dev/clir/middleware"
)

func TestFlags(t *testing.T) {
	t.Run("can set flags on a root route", func(t *testing.T) {
		r := clir.NewRouter()

		var v *bool
		r.Use(middleware.Flags(func(fs *flag.FlagSet) {
			v = fs.Bool("v", false, "")
		}))

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
		is.NotNil(t, v)
		is.True(t, *v)
	})

	t.Run("can set flags on the root and subroutes", func(t *testing.T) {
		r := clir.NewRouter()

		var v *bool
		r.Use(middleware.Flags(func(fs *flag.FlagSet) {
			v = fs.Bool("v", false, "")
		}))

		var called bool
		var fancy *bool

		r.Branch("dance", func(r *clir.Router) {
			r.Use(middleware.Flags(func(fs *flag.FlagSet) {
				fancy = fs.Bool("fancypants", false, "")
			}))

			r.RouteFunc("", func(ctx clir.Context) error {
				called = true
				return nil
			})
		})

		err := r.Run(clir.Context{
			Args: []string{"-v", "dance", "-fancypants"},
		})
		is.NotError(t, err)
		is.True(t, called)
		is.NotNil(t, v)
		is.True(t, *v)
		is.NotNil(t, fancy)
		is.True(t, *fancy)
	})

	t.Run("does not error on help flag but prints usage", func(t *testing.T) {
		for _, f := range []string{"-h", "-help"} {
			t.Run(f, func(t *testing.T) {
				r := clir.NewRouter()

				var v *bool
				var outerFS *flag.FlagSet
				r.Use(middleware.Flags(func(fs *flag.FlagSet) {
					outerFS = fs
					v = fs.Bool("v", false, "")
				}))

				var called bool
				r.RouteFunc("", func(ctx clir.Context) error {
					called = true
					return nil
				})

				var b strings.Builder
				err := r.Run(clir.Context{
					Args: []string{f},
					Err:  &b,
				})
				is.NotError(t, err)
				is.True(t, !called)
				is.NotNil(t, v)
				is.True(t, !*v)

				var usageB strings.Builder
				outerFS.SetOutput(&usageB)
				outerFS.Usage()
				is.Equal(t, usageB.String(), b.String())
			})
		}
	})
}

func ExampleFlags() {
	r := clir.NewRouter()

	var v *bool
	r.Use(middleware.Flags(func(fs *flag.FlagSet) {
		v = fs.Bool("v", false, "verbose output")
	}))

	r.RouteFunc("", func(ctx clir.Context) error {
		if *v {
			ctx.Println("Hello!")
		}
		return nil
	})

	_ = r.Run(clir.Context{
		Args: []string{"-v"},
		Out:  os.Stdout,
	})
	// Output: Hello!
}
