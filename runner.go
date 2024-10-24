// Package clir provides a definition of a [Router] and a [Runner].
package clir

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// Context for a [Runner] when it runs.
type Context struct {
	Args    []string
	Ctx     context.Context
	Err     io.Writer
	In      io.Reader
	Matches []string
	Out     io.Writer
}

func (c Context) Println(a ...any) {
	_, _ = fmt.Fprintln(c.Out, a...)
}

func (c Context) Printfln(format string, a ...any) {
	_, _ = fmt.Fprintf(c.Out, format+"\n", a...)
}

func (c Context) Errorln(a ...any) {
	_, _ = fmt.Fprintln(c.Err, a...)
}

func (c Context) Errorfln(format string, a ...any) {
	_, _ = fmt.Fprintf(c.Err, format+"\n", a...)
}

// Runner can run with a [Context].
type Runner interface {
	Run(ctx Context) error
}

// RunnerFunc is a function which satisfies [Runner].
type RunnerFunc func(ctx Context) error

// Run satisfies [Runner].
func (f RunnerFunc) Run(ctx Context) error {
	return f(ctx)
}

// Run a [Runner] with a default [Context], which is:
// - Get args from [os.Args]
// - Create context which is cancelled on [syscall.SIGTERM] or [syscall.SIGINT]
// - Use [os.Stdin] for input
// - Use [os.Stdout] for output
// - Use [os.Stderr] for errors
// - Prints to [os.Stderr] and calls os.Exit(1) on errors from [Runner.Run]
func Run(r Runner) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	runCtx := Context{
		Args: os.Args[1:],
		Ctx:  ctx,
		Err:  os.Stderr,
		In:   os.Stdin,
		Out:  os.Stdout,
	}

	if err := r.Run(runCtx); err != nil {
		runCtx.Errorln("Error:", err)
		os.Exit(1)
	}
}

var _ Runner = (*RunnerFunc)(nil)
