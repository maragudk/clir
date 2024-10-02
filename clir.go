package clir

import (
	"fmt"
	"io"
)

type Router struct {
	paths    []string
	commands map[string]Command
}

func NewRouter() *Router {
	return &Router{
		commands: map[string]Command{},
	}
}

type Context struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type Command = func(ctx Context) error

func (r *Router) Route(path string, cmd Command) {
	r.paths = append(r.paths, path)
	r.commands[path] = cmd
}

func Main(args []string, inR io.Reader, outW, errW io.Writer, r *Router) {
	if len(args) < 2 {

		return
	}

	path := args[1]
	if cmd, ok := r.commands[path]; ok {
		ctx := Context{
			In:  inR,
			Out: outW,
			Err: errW,
		}
		if err := cmd(ctx); err != nil {
			_, _ = fmt.Fprintln(errW, "Error:", err)
		}
	}
}
