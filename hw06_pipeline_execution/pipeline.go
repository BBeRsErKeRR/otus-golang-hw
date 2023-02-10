package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	worker := func(done In, stage Stage, inputC In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for r := range stage(inputC) {
				select {
				case <-done:
					return
				case out <- r:
				}
			}
		}()
		return out
	}

	reduce := func(done In, s []Stage, f func(done In, stage Stage, inputC In) Out, init In) Out {
		acc := init
		for _, v := range s {
			select {
			case <-done:
			default:
				acc = f(done, v, acc)
			}
		}
		return acc
	}

	return reduce(done, stages, worker, in)
}
