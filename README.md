# WorkerPool

[![GoDoc](https://godoc.org/github.com/qvik/workerpool?status.svg)](https://godoc.org/github.com/qvik/workerpool)

This is a worker pool implementation in Go language.

It provides a configurable concurrent task queue and a simple API.

## Installation

To install the library, run:

```sh
go get -u https://github.com/qvik/workerpool
```

## Usage

Here's an example of how to use the library:

```go
package main

import (
    "log"
    "github.com/qvik/workerpool"
)

func main() {
    // Create a WorkerPool with specified concurrency and queue size
    p := workerpool.NewWorkerPool(2, 2, 2)

    // It's a good strategy to process the results in a separate goroutine
    // instead of the goroutine that starts the tasks; this way the results
    // get consumed as soon as they become available and a possible
    // deadlock can be avoided.
    go func() {
        for i := 0; i < 2; i++ {
            res := <-p.GetResultsChannel()

            if res.Error != nil {
                // TODO: Handle error
            }
        }

        log.Printf("both tasks have completed")
    }()

    // Launch task 1
    p.AddTask(func() {
        // TODO: Do something here
    })

    // Launch task 2
    p.AddTask(func() {
        // TODO: Do something else here
    })

    // Close down the worker pool
    p.Close()
}
```

If you don't care about the results, you may simply call `p.WaitAll()` after all the tasks have been started to wait for their completion.

## License

The library is released under the [MIT license](LICENSE.md).

## Contact

Contact [Matti Dahlbom](mailto:matti@qvik.fi) if any questions arise.
