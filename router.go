package clir

// CommandRouter is a router for commands which itself satisfies [Command].
type CommandRouter struct {
	patterns []string
	commands map[string]Command
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		commands: map[string]Command{},
	}
}

// Run satisfies [Command].
func (c *CommandRouter) Run(ctx Context) error {
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

func (c *CommandRouter) Route(pattern string, cmd Command) {
	c.patterns = append(c.patterns, pattern)
	c.commands[pattern] = cmd
}

func (c *CommandRouter) RouteFunc(pattern string, cmd CommandFunc) {
	c.Route(pattern, cmd)
}

var _ Command = (*CommandRouter)(nil)
