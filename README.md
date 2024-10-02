# CLIR, the Command Line Interface Router

[![Go](https://github.com/maragudk/clir/actions/workflows/ci.yml/badge.svg)](https://github.com/maragudk/clir/actions/workflows/ci.yml)

You can think of command routing in a CLI the same way as routing in an HTTP server:
- Subcommands are URL paths
- Positional arguments are URL path parameters
- Flags are URL query parameters
- STDIN/STDOUT are the request/response bodies

CLIR is a Command Line Interface Router.

```shell
go get maragu.dev/clir
```

Made with ✨sparkles✨ by [maragu](https://www.maragu.dev/).
