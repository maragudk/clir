package clir

import (
	"fmt"
	"regexp"
	"strings"
)

// Router for [Runner]-s which itself satisfies [Runner].
type Router struct {
	middlewares []Middleware
	patterns    []*regexp.Regexp
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
	// Apply middlewares first, because they can modify the context, including the Context.Args to match against.
	var middlewareCtx Context
	var runner Runner = RunnerFunc(func(ctx Context) error {
		middlewareCtx = ctx
		return nil
	})
	// Apply middlewares in reverse order, so the first middleware is the outermost one, to be called first.
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		runner = r.middlewares[i](runner)
	}
	if err := runner.Run(ctx); err != nil {
		return fmt.Errorf("error while applying middleware: %w", err)
	}
	ctx = middlewareCtx

	for _, pattern := range r.patterns {
		if (len(ctx.Args) == 0 && pattern.String() == "^$") || (len(ctx.Args) > 0 && pattern.MatchString(ctx.Args[0])) {
			runner = r.runners[pattern.String()]
			if len(ctx.Args) > 0 {
				ctx.Matches = pattern.FindStringSubmatch(ctx.Args[0])
				ctx.Args = ctx.Args[1:]
			}

			return runner.Run(ctx)
		}
	}

	//for _, router := range r.routers {
	//	if err := router.Run(ctx); err == nil {
	//		return err
	//	}
	//}

	return ErrorRouteNotFound
}

// Route a [Runner] with the given pattern.
// Routes are matched in the order they were added.
func (r *Router) Route(pattern string, runner Runner) {
	if !strings.HasPrefix(pattern, "^") {
		pattern = "^" + pattern
	}
	if !strings.HasSuffix(pattern, "$") {
		pattern += "$"
	}

	if _, ok := r.runners[pattern]; ok {
		panic("cannot add route which already exists")
	}

	r.patterns = append(r.patterns, regexp.MustCompile(pattern))
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
// The middlewares from the parent router are used in the new router,
// but new middlewares within the scope are only added to the new router, not the parent router.
func (r *Router) Scope(cb func(r *Router)) {
	panic("not implemented")
	//newR := NewRouter()
	//cb(newR)
	//r.routers = append(r.routers, newR)
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
