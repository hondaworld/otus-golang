package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	wg := sync.WaitGroup{}
	//out := make(Bi)
	//mu := sync.Mutex{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		//defer close(out)

		//isFinished := false

		for _, stage := range stages {
			select {
			case <-done:
				//isFinished = true
				println(111)
				return
			default:
				in = stage(in)
			}
		}

		//if !isFinished {
		//	out <- in
		//}

	}()

	wg.Wait()

	return in
	//close(out)

	//if isFinished {
	//	return nil
	//}

}
