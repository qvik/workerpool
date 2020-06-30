# WorkerPool

[![GoDoc](https://godoc.org/github.com/qvik/qvik-go-workerpool?status.svg)](https://godoc.org/github.com/qvik/qvik-go-workerpool)

This is a worker pool implementation in Go language.

It provides a configurable concurrent task queue and a simple API.

## Installation

To install the library, run:

```sh
go get -u https://github.com/qvik/qvik-go-workerpool
```

## Usage

Here's an example of how to use the library:

```go
package main

import (
    "github.com/qvik/workerpool"
)

func main() {
    // Create a WorkerPool with specified concurrency and queue size
    p := workerpool.NewWorkerPool(2, 2)

    p.AddTask(func() {
        // TODO: Do something here
    })

    p.AddTask(func() {
        // TODO: Do something else here
    })

    // Wait till both tasks have completed
    p.WaitAll()

    // Close down the worker pool
    p.Close()
}
```

## License

The library is released under the [MIT license](LICENSE.md).

## Contact

Contact [Matti Dahlbom](mailto:matti@qvik.fi) if any questions arise.
