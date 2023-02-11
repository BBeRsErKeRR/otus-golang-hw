package hw06pipelineexecution

import (
	"time"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	worker := func(done In, stageReader In) Out {
		// limiter to send data to receive done signal
		limiter := time.NewTicker(5 * time.Millisecond)
		// result chan
		out := make(Bi)
		go func() {
			defer close(out)
			for {
				// check done signal before read stage out
				select {
				case <-done:
					return
				case stOut, ok := <-stageReader:
					if !ok {
						return
					}
					// wait tick and check done signal
					<-limiter.C
					select {
					case <-done:
						return
					case out <- stOut:
					}
				}
			}
		}()
		return out
	}

	// concatenate all results
	reduce := func(done In, s []Stage, f func(done In, stageReader In) Out, init In) Out {
		acc := init
		for _, v := range s {
			acc = f(done, v(acc))
		}
		return acc
	}

	return reduce(done, stages, worker, in)
}
