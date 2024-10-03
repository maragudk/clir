package clir_test

import (
	"strings"
	"testing"

	"maragu.dev/is"

	"maragu.dev/clir"
)

func TestRun(t *testing.T) {
	t.Run("can run a command", func(t *testing.T) {
		err := clir.Run(clir.CommandFunc(func(ctx clir.Context) error {
			is.True(t, strings.Contains(ctx.Args[0], "clir.test"))
			return nil
		}))
		is.NotError(t, err)
	})
}
