package clir

// CommandMux is a multiplexer/router for commands which itself satisfies [Command].
type CommandMux struct {
	patterns []string
	commands map[string]Command
}

func NewCommandMux() *CommandMux {
	return &CommandMux{
		commands: map[string]Command{},
	}
}

func (c *CommandMux) Run(ctx Context) error {
	if len(ctx.Args) == 0 {
		root, ok := c.commands[""]
		if !ok {
			return ErrorNotFound
		}
		return root.Run(ctx)
	}

	for _, pattern := range c.patterns {
		if ctx.Args[0] == pattern {
			cmd := c.commands[pattern]
			ctx.Args = ctx.Args[1:]
			return cmd.Run(ctx)
		}
	}

	return ErrorNotFound
}

func (c *CommandMux) Handle(pattern string, cmd Command) {
	c.patterns = append(c.patterns, pattern)
	c.commands[pattern] = cmd
}

func (c *CommandMux) HandleFunc(pattern string, cmd CommandFunc) {
	c.Handle(pattern, cmd)
}

var _ Command = (*CommandMux)(nil)
