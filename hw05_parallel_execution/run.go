package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m < 0 {
		return ErrErrorsLimitExceeded
	}
	var errorsC int32
	wg := sync.WaitGroup{}
	ch := make(chan Task, n)
	wg.Add(n + 1)

	// Producer
	go func() {
		defer wg.Done()
		defer close(ch)
		for _, task := range tasks {
			// close chan if found errors
			if int(atomic.LoadInt32(&errorsC)) >= m {
				break
			}
			ch <- task
		}
	}()

	// Consumer's
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for f := range ch {
				if int(atomic.LoadInt32(&errorsC)) > m {
					break
				}
				err := f()
				if err != nil {
					atomic.AddInt32(&errorsC, 1)
				}
			}
		}()
	}

	wg.Wait()

	// Check last errors
	if int(errorsC) > m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
