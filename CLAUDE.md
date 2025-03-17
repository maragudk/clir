# CLAUDE.md - Clir Repository Guidelines

## Build/Test Commands
- Build: `go build ./...`
- Test all: `make test` (runs with coverage and shuffle)
- Test single: `go test ./path/to/package -run "TestName"`
- Run specific subtest: `go test ./middleware -run "^TestFlags$/does_not_error_on_help_flag"`
- Lint: `make lint` (uses golangci-lint)
- Coverage: `make cover` (view HTML report)
- Benchmark: `make benchmark`

## Code Style Guidelines
- **Package**: Main package is `clir`, with `middleware` subpackage
- **Imports**: Alphabetized; standard lib first, blank line, then third-party
- **Naming**: CamelCase for exported items; lowercase for internal
- **Error Handling**: Custom error types with constants for common errors
- **Tests**: Table-driven with `t.Run`, use `maragu.dev/is` package for assertions
- **Documentation**: Complete doc comments on all exported functions
- **Design Patterns**:
  - Middleware functions return `clir.Middleware` type
  - Router supports nested routes with branch functionality
  - Return errors rather than panicking (except during setup)
