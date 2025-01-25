package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type StageOut func(in In, done In) (out Out)

func pipeline(stages ...Stage) StageOut {
	return func(in In, done In) (out Out) {
		for _, stage := range stages {
			in = take(in, done, stage)
		}

		return in
	}
}

func take(in In, done In, stage Stage) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for v := range stage(in) {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	finalStage := pipeline(stages...)

	go func() {
		defer close(out)

		for v := range finalStage(in, done) {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}
