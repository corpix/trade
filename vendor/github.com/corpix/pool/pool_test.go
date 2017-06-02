package pool

// The MIT License (MIT)
//
// Copyright © 2017 Dmitry Moskowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoolParallel(t *testing.T) {
	tasks := 10
	workers := 10
	queue := 10
	sleep := 100 * time.Millisecond

	p := New(workers, queue)
	defer p.Close()

	w := &sync.WaitGroup{}
	w.Add(tasks)

	started := time.Now()
	for n := 0; n < tasks; n++ {
		p.Feed <- NewWork(
			context.Background(),
			func(ctx context.Context) {
				select {
				case <-ctx.Done():
				case <-time.After(sleep):
				}
				w.Done()
			},
		)
	}
	w.Wait()
	finished := time.Now()

	assert.False(
		t,
		started.Add(sleep*time.Duration(tasks)).Before(finished),
	)
}

func TestPoolContextCancel(t *testing.T) {
	tasks := 5
	workers := 5
	queue := 0

	p := New(workers, queue)
	defer p.Close()

	w := &sync.WaitGroup{}
	w.Add(tasks)

	cancels := make(chan int, tasks*2)
	defer close(cancels)

	for n := 0; n < tasks; n++ {
		ctx, cancel := context.WithCancel(
			context.Background(),
		)
		cancel()

		p.Feed <- NewWork(
			ctx,
			func(ctx context.Context) {
				select {
				case <-ctx.Done():
					cancels <- 1
				}
				w.Done()
			},
		)
	}

	canceled := 0
	go func() {
		for c := range cancels {
			canceled += c
			if canceled == tasks {
				break
			}
		}
		w.Done()
	}()
	w.Add(1)

	w.Wait()

	assert.Equal(t, tasks, canceled)
}
