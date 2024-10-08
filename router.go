package clir

// CommandRouter is a router for commands which itself satisfies [Command].
type CommandRouter struct {
	commands    map[string]Command
	middlewares []Middleware
	patterns    []string
	subRouters  []*CommandRouter
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		commands: map[string]Command{},
	}
}

// Run satisfies [Command].
func (c *CommandRouter) Run(ctx Context) error {
	if len(ctx.Args) == 0 {
		cmd, ok := c.commands[""]
		if !ok {
			return ErrorNotFound
		}

		// Apply middlewares in reverse order, so that middlewares are applied in the order they were added.
		for i := len(c.middlewares) - 1; i >= 0; i-- {
			cmd = c.middlewares[i](cmd)
		}

		return cmd.Run(ctx)
	}

	for _, pattern := range c.patterns {
		if ctx.Args[0] == pattern {
			cmd := c.commands[pattern]
			ctx.Args = ctx.Args[1:]

			for i := len(c.middlewares) - 1; i >= 0; i-- {
				cmd = c.middlewares[i](cmd)
			}

			return cmd.Run(ctx)
		}
	}

	for _, r := range c.subRouters {
		if err := r.Run(ctx); err == nil {
			return err
		}
	}

	return ErrorNotFound
}

// Sub adds a subcommand to the router with the given pattern.
func (c *CommandRouter) Sub(pattern string, cmd Command) {
	c.patterns = append(c.patterns, pattern)
	c.commands[pattern] = cmd
}

// SubFunc is a convenience method for adding a subcommand with a [CommandFunc].
// It's the same as calling [CommandRouter.Sub] with a [CommandFunc], but you don't have to wrap the function.
func (c *CommandRouter) SubFunc(pattern string, cmd CommandFunc) {
	c.Sub(pattern, cmd)
}

// Group commands with a new router.
// The middleware from the parent router is copied to the new router.
func (c *CommandRouter) Group(cb func(r *CommandRouter)) {
	r := NewCommandRouter()
	r.middlewares = append(r.middlewares, c.middlewares...)
	cb(r)
	c.subRouters = append(c.subRouters, r)
}

type Middleware = func(next Command) Command

// Use a middleware for all commands. If called on the root router, it will apply to all commands.
// If called in a Group, it will apply to all commands in that group.
func (c *CommandRouter) Use(middlewares ...Middleware) {
	if len(c.commands) > 0 {
		panic("cannot add middlewares after adding commands")
	}
	c.middlewares = append(c.middlewares, middlewares...)
}

var _ Command = (*CommandRouter)(nil)
