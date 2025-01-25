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
	temp := make(Bi)

	go func() {
		defer close(out)

		for v := range temp {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	go func() {
		defer close(temp)

		for v := range stage(in) {
			temp <- v
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	temp := make(Bi)

	finalStage := pipeline(stages...)

	go func() {
		defer close(out)

		for v := range temp {
			select {
			case <-done:
				close(temp)
				return
			case out <- v:
			}
		}
	}()

	go func() {
		defer close(temp)

		for v := range finalStage(in, done) {
			temp <- v
		}
	}()

	return out
}
