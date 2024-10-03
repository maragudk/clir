// Package clir provides a definition of a runnable [Command] as well as a [CommandMux], which is a multiplexer/router for commands.
package clir

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// Context for a [Command] that runs.
type Context struct {
	Args []string
	Ctx  context.Context
	Err  io.Writer
	In   io.Reader
	Out  io.Writer
}

// Command can be run with a Context.
type Command interface {
	Run(ctx Context) error
}

// CommandFunc is a function which satisfies [Command].
type CommandFunc func(ctx Context) error

// Run satisfies [Command].
func (f CommandFunc) Run(ctx Context) error {
	return f(ctx)
}

// Run a [Command] with default options, which are:
// - Get args from [os.Args]
// - Create context which is cancelled on [syscall.SIGTERM] or [syscall.SIGINT]
// - Use [os.Stdin] for input
// - Use [os.Stdout] for output
// - Use [os.Stderr] for errors
func Run(cmd Command) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cmdCtx := Context{
		Args: os.Args,
		Ctx:  ctx,
		Err:  os.Stderr,
		In:   os.Stdin,
		Out:  os.Stdout,
	}

	return cmd.Run(cmdCtx)
}

var _ Command = (*CommandFunc)(nil)
