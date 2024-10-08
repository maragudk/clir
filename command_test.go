package clir_test

import (
	"testing"

	"maragu.dev/is"

	"maragu.dev/clir"
)

func TestRun(t *testing.T) {
	t.Run("can run a command", func(t *testing.T) {
		var called bool
		clir.Run(clir.CommandFunc(func(ctx clir.Context) error {
			called = true
			return nil
		}))
		is.True(t, called)
	})
}
