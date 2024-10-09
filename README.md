# CLIR, the Command Line Interface Router

<img src="logo.jpg" alt="CLIR logo" width="300" align="right"/>

[![GoDoc](https://pkg.go.dev/badge/maragu.dev/clir)](https://pkg.go.dev/maragu.dev/clir)
[![Go](https://github.com/maragudk/clir/actions/workflows/ci.yml/badge.svg)](https://github.com/maragudk/clir/actions/workflows/ci.yml)

You can think of routing in a CLI the same way as routing in an HTTP server:
- Subcommands are URL paths
- Positional arguments are URL path parameters
- Flags are URL query parameters
- STDIN/STDOUT are the request/response bodies

CLIR is a Command Line Interface Router. ⚠️ Work-in-progress and will probably change a lot, but play around with it!

Made with ✨sparkles✨ by [maragu](https://www.maragu.dev/).

Does your company depend on this project? [Contact me at markus@maragu.dk](mailto:markus@maragu.dk?Subject=Supporting%20your%20project) to discuss options for a one-time or recurring invoice to ensure its continued thriving.

## Usage

```shell
go get maragu.dev/clir
```

```go
package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"maragu.dev/clir"
)

func main() {
	// Initialize dependencies
	l := slog.New(slog.NewTextHandler(os.Stderr, nil))
	c := &http.Client{
		Timeout: time.Second,
	}

	// Create a new router which is also something that can be run.
	r := clir.NewRouter()

	// Add logging middleware to all routes.
	r.Use(log(l))

	// Add a root route which calls printHello.
	r.Route("", printHello())

	// Scope some middleware to just the routes within the scope.
	r.Scope(func(r *clir.Router) {
		r.Use(ping(c))

		r.Route("get", get(c))
	})

	// Branch with subcommands
	r.Branch("post", func(r *clir.Router) {
		r.Use(ping(c))

		r.Route("stdin", postFromStdin(c))
		r.Route("random", postFromRandom(c))
	})

	// Run the router with a default clir.Context.
	clir.Run(r)
}

// printHello to stdout.
func printHello() clir.RunnerFunc {
	return func(ctx clir.Context) error {
		ctx.Println("Hello!")

		return nil
	}
}

// get example.com.
func get(c *http.Client) clir.RunnerFunc {
	return func(ctx clir.Context) error {
		res, err := c.Get("https://example.com")
		if err != nil {
			ctx.Errorln("Didn't get it.")
			return err
		}

		ctx.Println("Got it! Response:", res.Status)
		return nil
	}
}

// postFromStdin to example.com.
func postFromStdin(c *http.Client) clir.RunnerFunc {
	return func(ctx clir.Context) error {
		res, err := c.Post("https://example.com", "text/plain", ctx.In)
		if err != nil {
			ctx.Errorln("Didn't post stdin.")
			return err
		}

		ctx.Println("Posted stdin! Response:", res.Status)
		return nil
	}
}

func postFromRandom(c *http.Client) clir.RunnerFunc {
	return func(ctx clir.Context) error {
		randomNumber := rand.Int()
		ctx.Println("Random number is", randomNumber)

		res, err := c.Post("https://example.com", "text/plain", strings.NewReader(fmt.Sprint(randomNumber)))
		if err != nil {
			ctx.Errorln("Didn't post the random number.")
			return err
		}

		ctx.Println("Posted the random number! Response:", res.Status)
		return nil
	}
}

// log the arguments to the given [slog.Logger].
func log(l *slog.Logger) clir.Middleware {
	return func(next clir.Runner) clir.Runner {
		return clir.RunnerFunc(func(ctx clir.Context) error {
			l.InfoContext(ctx.Ctx, "Called", "args", ctx.Args)

			// Remember to call next.Run, or the chain will stop here.
			return next.Run(ctx)
		})
	}
}

// ping a URL to check the network.
func ping(c *http.Client) clir.Middleware {
	return func(next clir.Runner) clir.Runner {
		return clir.RunnerFunc(func(ctx clir.Context) error {
			if _, err := c.Get("https://example.com"); err != nil {
				return err
			}
			return next.Run(ctx)
		})
	}
}
```
