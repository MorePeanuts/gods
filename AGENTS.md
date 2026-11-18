# Repository Guidelines

## Project Structure

`gods` is a generic Go data structure library. The module is `github.com/MorePeanuts/gods`. Shared contracts live in `container/` and in the top-level container packages. Concrete implementations include `list/arraylist`, `list/linkedlist`, `maps/treemap`, and `set/treeset`. Reusable structures live in `tree/rbtree` and `tree/binaryheap`. The tree map and priority queue wrap those implementations. Put tests next to the code they exercise in `*_test.go` files.

## Build, Test, and Documentation

Run commands from the repository root with the Go version declared in `go.mod`.

- `go build ./...` compiles every package.
- `go test ./...` runs the full suite. Use `go test ./tree/...` to focus on tree packages.
- `go vet ./...` checks for common Go mistakes.
- `go doc -http` opens local API documentation in a browser. Use `go doc ./queue/priorityqueue` for a terminal view.

## Coding Style

Run `gofmt -w` on changed Go files and use standard Go naming and indentation. Write English GoDoc for exported packages, types, and methods. Keep balancing and heap maintenance in `tree/` when a public container can delegate to an existing implementation. Preserve documented empty-container behavior, comparator ordering, and shallow-copy semantics. Document any change to a public contract.

## Testing

Use Go's standard `testing` package and name tests `TestXxx`. For behavior changes, exercise the public API and relevant edge cases, such as empty containers, duplicate keys, custom comparators, or independent clone storage. For balancing changes, check the structure's invariants as well as returned values.

## Commits and Pull Requests

Use English Conventional Commit subjects, such as `feat: add binary heap and priority queue` or `refactor: extract reusable red-black tree`. Mark breaking API changes with `!` and explain the migration in a `BREAKING CHANGE:` footer. In pull requests, summarize the behavior change, list verification commands, and link the relevant issue when available.
