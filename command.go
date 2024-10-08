// Package clir provides a definition of a runnable [Command] as well as a [CommandRouter], which is a multiplexer/router for commands.
package clir

import (
	"context"
	"fmt"
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

// Command can be run with a [Context].
type Command interface {
	Run(ctx Context) error
}

// CommandFunc is a function which satisfies [Command].
type CommandFunc func(ctx Context) error

// Run satisfies [Command].
func (f CommandFunc) Run(ctx Context) error {
	return f(ctx)
}

// Run a [Command] with default a [Context], which is:
// - Get args from [os.Args]
// - Create context which is cancelled on [syscall.SIGTERM] or [syscall.SIGINT]
// - Use [os.Stdin] for input
// - Use [os.Stdout] for output
// - Use [os.Stderr] for errors
// - Prints to [os.Stderr] and calls os.Exit(1) on errors from the command
func Run(cmd Command) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cmdCtx := Context{
		Args: os.Args[1:],
		Ctx:  ctx,
		Err:  os.Stderr,
		In:   os.Stdin,
		Out:  os.Stdout,
	}

	if err := cmd.Run(cmdCtx); err != nil {
		cmdCtx.Errorln("Error:", err)
		os.Exit(1)
	}
}

var _ Command = (*CommandFunc)(nil)
