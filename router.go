package clir

// Router for [Runner]-s which itself satisfies [Runner].
type Router struct {
	middlewares []Middleware
	patterns    []string
	routers     []*Router
	runners     map[string]Runner
}

func NewRouter() *Router {
	return &Router{
		runners: map[string]Runner{},
	}
}

// Run satisfies [Runner].
func (r *Router) Run(ctx Context) error {
	for _, pattern := range r.patterns {
		if (len(ctx.Args) == 0 && pattern == "") || (ctx.Args[0] == pattern) {
			runner := r.runners[pattern]
			if len(ctx.Args) > 0 {
				ctx.Args = ctx.Args[1:]
			}

			for i := len(r.middlewares) - 1; i >= 0; i-- {
				runner = r.middlewares[i](runner)
			}

			return runner.Run(ctx)
		}
	}

	for _, router := range r.routers {
		if err := router.Run(ctx); err == nil {
			return err
		}
	}

	return ErrorRouteNotFound
}

// Route a [Runner] with the given pattern.
// Routes are matched in the order they were added.
func (r *Router) Route(pattern string, runner Runner) {
	if _, ok := r.runners[pattern]; ok {
		panic("cannot add route which already exists")
	}
	r.patterns = append(r.patterns, pattern)
	r.runners[pattern] = runner
}

// RouteFunc is like [Router.Route], but with a [RunnerFunc].
func (r *Router) RouteFunc(pattern string, runner RunnerFunc) {
	r.Route(pattern, runner)
}

// Branch into a new [Router] with the given pattern.
func (r *Router) Branch(pattern string, cb func(r *Router)) {
	newR := NewRouter()
	cb(newR)
	r.Route(pattern, newR)
}

// Scope into a new [Router].
// The middlewares from the parent router are copied to the new router,
// but new middlewares within the scope are only added to the new router, not the parent router.
func (r *Router) Scope(cb func(r *Router)) {
	newR := NewRouter()
	newR.middlewares = append(newR.middlewares, r.middlewares...)
	cb(newR)
	r.routers = append(r.routers, newR)
}

// Middleware for [Router.Use].
type Middleware = func(next Runner) Runner

// Use [Middleware] on the current branch of the [Router].
// If called in a [Scope], it will apply to all routes in that scope.
func (r *Router) Use(middlewares ...Middleware) {
	if len(r.runners) > 0 {
		panic("cannot add middlewares after adding routes")
	}
	r.middlewares = append(r.middlewares, middlewares...)
}

var _ Runner = (*Router)(nil)
