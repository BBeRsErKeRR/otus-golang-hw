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
	var errorsC int32
	var lenT int32
	if int(errorsC) > m {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	ch := make(chan Task, n)
	wg.Add(n + 1)

	// Producer
	go func(sendTasksCountAddr *int32, errorsCountAddr *int32) {
		defer wg.Done()
		tasksLen := len(tasks)
		for {
			// close chan if found errors
			testE := int(atomic.LoadInt32(errorsCountAddr)) >= m
			if testE {
				close(ch)
				break
			}

			// Add new items into chan
			sendTasksCount := int(atomic.LoadInt32(sendTasksCountAddr))
			test := sendTasksCount < tasksLen
			if test {
				ch <- tasks[sendTasksCount]
				atomic.AddInt32(sendTasksCountAddr, 1)
			} else {
				close(ch)
				break
			}
		}
	}(&lenT, &errorsC)

	// Consumer's
	for i := 0; i < n; i++ {
		go func(errorsCountAddr *int32) {
			defer wg.Done()
			for {
				test := int(atomic.LoadInt32(errorsCountAddr)) <= m
				if test {
					f, ok := <-ch
					if !ok {
						break
					}
					err := f()
					if err != nil {
						atomic.AddInt32(errorsCountAddr, 1)
					}
				} else {
					break
				}
			}
		}(&errorsC)
	}

	wg.Wait()

	// Check last errors
	if int(errorsC) > m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
