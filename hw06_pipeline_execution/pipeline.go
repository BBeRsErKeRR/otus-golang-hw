package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Resolve stage worker
	worker := func(stage Stage, input In, done In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for r := range stage(input) {
				select {
				case <-done:
					return
				default:
					out <- r
				}
			}
		}()
		return out
	}
	// Buffer to reduce stages results into one value
	// TODO: create reducer function
	buffer := make(map[int](Out))
	for i, stage := range stages {
		if i == 0 {
			buffer[i] = worker(stage, in, done)
		} else {
			buffer[i] = worker(stage, buffer[i-1], done)
		}
	}
	return buffer[len(stages)-1]
}
