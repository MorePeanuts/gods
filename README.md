# gods

A generic Go data structure library with lists, stacks, queues, maps, sets, and reusable tree implementations.

The original GoDS project inspired this library, which we started because upstream maintenance appeared to have stopped. [emirpasic/gods](https://github.com/emirpasic/gods)

## Packages

| Category | Implementations |
| --- | --- |
| Lists | [Array list](list/arraylist), [doubly linked list](list/linkedlist) |
| Queues | [FIFO queue](queue), [deque](deque), [priority queue](queue/priorityqueue) |
| Stacks | [Slice-backed stack](stack) |
| Maps | [Hash map](maps/hashmap), [insertion-ordered hash map](maps/linkedhashmap), [tree map](maps/treemap) |
| Sets | [Hash set](set/hashset), [insertion-ordered hash set](set/linkedhashset), [tree set](set/treeset) |
| Trees | [Red-black tree](tree/rbtree), [binary heap](tree/binaryheap) |

## Quick start

Add the module to a Go project.

```sh
go get github.com/MorePeanuts/gods
```

The default priority queue dequeues the smallest element first. Use `NewWith` to supply a comparator for a different priority order.

```go
package main

import (
	"fmt"

	"github.com/MorePeanuts/gods/queue/priorityqueue"
)

func main() {
	q := priorityqueue.New(3, 1, 2)
	for !q.IsEmpty() {
		value, _ := q.Dequeue()
		fmt.Println(value) // Prints 1, 2, and 3 on separate lines.
	}
}
```

## Development

Use the Go version declared in [go.mod](go.mod). Run these commands from the repository root.

```sh
make build  # Compile all packages.
make fmt    # Format Go files with gofmt.
make check  # Check formatting, run go vet, and run tests.
make doc    # Open local API documentation in a browser.
```

Run `make help` for all targets. See [AGENTS.md](AGENTS.md) for contributor conventions.
