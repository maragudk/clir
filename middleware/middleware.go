// Package middleware provides useful middleware for a [clir.Router].
package middleware

import (
	"errors"
	"flag"

	"maragu.dev/clir"
)

// Flags middleware allows you to set flags on a route.
func Flags(cb func(fs *flag.FlagSet)) clir.Middleware {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	cb(fs)

	return func(next clir.Runner) clir.Runner {
		return clir.RunnerFunc(func(ctx clir.Context) error {
			fs.SetOutput(ctx.Err)
			if err := fs.Parse(ctx.Args); err != nil {
				if errors.Is(err, flag.ErrHelp) {
					return nil
				}
				return err
			}
			ctx.Args = fs.Args()
			return next.Run(ctx)
		})
	}
}
