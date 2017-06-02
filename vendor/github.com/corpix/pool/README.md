pool
----

[![Build Status](https://travis-ci.org/corpix/pool.svg?branch=master)](https://travis-ci.org/corpix/pool)

Simplest goroutine pool ever.

## Simple example
``` go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/corpix/pool"
)

func main() {
	p := pool.New(10, 10)
	defer p.Close()

	w := &sync.WaitGroup{}

	tasks := 10
	sleep := 1 * time.Second

	for n := 0; n < tasks; n++ {
		w.Add(1)
		p.Feed <- pool.NewWork(
			context.Background(),
			func(n int) pool.Executor {
				return func(ctx context.Context) {
					select {
					case <-ctx.Done():
					default:
						time.Sleep(sleep)
						fmt.Printf("Finished work '%d'\n", n)
					}
					w.Done()
				}
			}(n),
		)
	}

	w.Wait()
}
```

Output:

> Results may differ on your machine, order is not guarantied.

``` console
$ go run ./example/simple/simple.go
Finished work '6'
Finished work '9'
Finished work '7'
Finished work '5'
Finished work '4'
Finished work '8'
Finished work '2'
Finished work '0'
Finished work '3'
Finished work '1'
```

## Example with work result

``` go
package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/corpix/pool"
)

func main() {
	p := pool.New(10, 10)
	defer p.Close()

	w := &sync.WaitGroup{}

	tasks := 10
	results := make(chan *pool.Result)
	defer close(results)

	for n := 0; n < tasks; n++ {
		w.Add(1)
		p.Feed <- pool.NewWorkWithResult(
			context.Background(),
			results,
			func(n int) pool.ExecutorWithResult {
				return func(ctx context.Context) (interface{}, error) {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					default:
						fmt.Printf("Finished work '%d'\n", n)
					}
					return n, nil
				}
			}(n),
		)
	}

	go func() {
		// Releasing one worker per iteration
		// when using not buffered channels.
		for result := range results {
			fmt.Printf(
				"Received result '%d'\n",
				result.Value.(int),
			)
			w.Done()
		}
	}()

	w.Wait()
}
```

Output:

> Results may differ on your machine, order is not guarantied.


``` console
$ go run ./example/with_result/with_result.go
Finished work '0'
Finished work '1'
Received result '0'
Received result '1'
Finished work '2'
Finished work '3'
Received result '2'
Received result '3'
Finished work '4'
Finished work '5'
Received result '4'
Received result '5'
Finished work '6'
Finished work '8'
Received result '6'
Received result '8'
Finished work '9'
Received result '9'
Finished work '7'
Received result '7'
```

## License

MIT
