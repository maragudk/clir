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

func ExampleArgs() {
	r := clir.NewRouter()

	var name *string
	r.Use(middleware.Args(func(as *middleware.ArgSet) {
		name = as.String("name", "world", "name to greet")
	}))

	r.RouteFunc("", func(ctx clir.Context) error {
		ctx.Printfln("Hello, %s!", *name)
		return nil
	})

	_ = r.Run(clir.Context{
		Args: []string{"Funky Person"},
		Out:  os.Stdout,
	})
	// Output: Hello, Funky Person!
}

func TestArgs(t *testing.T) {
	t.Run("can set args on a root route", func(t *testing.T) {
		r := clir.NewRouter()

		var name *string
		r.Use(middleware.Args(func(as *middleware.ArgSet) {
			name = as.String("name", "default", "set a name")
		}))

		var called bool
		r.RouteFunc("", func(ctx clir.Context) error {
			called = true
			return nil
		})

		err := r.Run(clir.Context{
			Args: []string{"fancypants"},
		})
		is.NotError(t, err)
		is.True(t, called)
		is.NotNil(t, name)
		is.Equal(t, "fancypants", *name)
	})

	t.Run("can set multiple args in order", func(t *testing.T) {
		r := clir.NewRouter()

		var name *string
		var age *string
		r.Use(middleware.Args(func(as *middleware.ArgSet) {
			name = as.String("name", "default", "set a name")
			age = as.String("age", "0", "set an age")
		}))

		var called bool
		r.RouteFunc("", func(ctx clir.Context) error {
			called = true
			return nil
		})

		err := r.Run(clir.Context{
			Args: []string{"alice", "30"},
		})
		is.NotError(t, err)
		is.True(t, called)
		is.NotNil(t, name)
		is.Equal(t, "alice", *name)
		is.NotNil(t, age)
		is.Equal(t, "30", *age)
	})

	t.Run("falls back to default if arg is missing", func(t *testing.T) {
		r := clir.NewRouter()

		var name *string
		var age *string
		r.Use(middleware.Args(func(as *middleware.ArgSet) {
			name = as.String("name", "default", "set a name")
			age = as.String("age", "0", "set an age")
		}))

		var called bool
		r.RouteFunc("", func(ctx clir.Context) error {
			called = true
			return nil
		})

		err := r.Run(clir.Context{
			Args: []string{"alice"},
		})
		is.NotError(t, err)
		is.True(t, called)
		is.NotNil(t, name)
		is.Equal(t, "alice", *name)
		is.NotNil(t, age)
		is.Equal(t, "0", *age) // Default value
	})

	t.Run("supports all positional argument types", func(t *testing.T) {
		// Setup a test to verify the different positional argument types
		var stringVal *string
		var intVal *int
		var boolVal *bool
		var floatVal *float64

		middleware := middleware.Args(func(as *middleware.ArgSet) {
			stringVal = as.String("string", "default", "string positional arg")
			intVal = as.Int("int", 0, "int positional arg")
			boolVal = as.Bool("bool", false, "bool positional arg")
			floatVal = as.Float64("float", 0.0, "float positional arg")
		})

		runner := middleware(clir.RunnerFunc(func(ctx clir.Context) error {
			return nil
		}))

		// Run with arguments matching our positional arguments
		err := runner.Run(clir.Context{
			Args: []string{"hello", "42", "true", "3.14"},
		})

		is.NotError(t, err)
		is.NotNil(t, stringVal)
		is.Equal(t, "hello", *stringVal)
		is.NotNil(t, intVal)
		is.Equal(t, 42, *intVal)
		is.NotNil(t, boolVal)
		is.Equal(t, true, *boolVal)
		is.NotNil(t, floatVal)
		is.Equal(t, 3.14, *floatVal)
	})

	t.Run("handles remaining args correctly", func(t *testing.T) {
		// Setup a test to verify our Args implementation passes the remaining arguments correctly
		var argValues []string

		middleware := middleware.Args(func(as *middleware.ArgSet) {
			as.String("first", "default", "first positional arg")
		})

		runner := middleware(clir.RunnerFunc(func(ctx clir.Context) error {
			argValues = ctx.Args
			return nil
		}))

		// Run with three arguments, first should be consumed by the positional arg, remaining should be passed through
		err := runner.Run(clir.Context{
			Args: []string{"value1", "extra1", "extra2"},
		})

		is.NotError(t, err)
		is.Equal(t, 2, len(argValues))
		is.Equal(t, "extra1", argValues[0])
		is.Equal(t, "extra2", argValues[1])
	})

	t.Run("works with subroutes", func(t *testing.T) {
		r := clir.NewRouter()

		var outerCalled, innerCalled bool
		r.RouteFunc("", func(ctx clir.Context) error {
			outerCalled = true
			return nil
		})

		var command *string
		r.Branch("run", func(r *clir.Router) {
			r.Use(middleware.Args(func(as *middleware.ArgSet) {
				command = as.String("command", "", "command to run")
			}))

			r.RouteFunc("", func(ctx clir.Context) error {
				innerCalled = true
				return nil
			})
		})

		err := r.Run(clir.Context{
			Args: []string{"run", "job"},
		})
		is.NotError(t, err)
		is.True(t, !outerCalled)
		is.True(t, innerCalled)
		is.NotNil(t, command)
		is.Equal(t, "job", *command)
	})
}
