package hw06pipelineexecution

import "fmt"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func concurrentPipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		select {
		case <-done:
			return done
		default:
			in = stage(in)
		}
	}

	return in
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	//wg := sync.WaitGroup{}
	//
	//wg.Add(1)

	pipeline := concurrentPipeline(in, done, stages...)

	go func() {
		defer close(out)

		for stage := range pipeline {
			select {
			case <-done:
				fmt.Println(111)
				return
			default:
				out <- stage
			}
		}
	}()

	//wg.Wait()

	return out

	//for value := range in {
	//	wg := sync.WaitGroup{}
	//	wg.Add(1)
	//	out := make(Bi)
	//
	//	go func() {
	//		defer wg.Done()
	//		defer close(out)
	//
	//		select {
	//		case <-done:
	//			return
	//		default:
	//			fmt.Println(1222)
	//			out <- value
	//			fmt.Println(3333)
	//			for _, stage := range stages {
	//				select {
	//				case <-done:
	//					//isFinished = true
	//					println(111)
	//					return
	//				case out <- stage(out):
	//				}
	//			}
	//		}
	//
	//		//if !isFinished {
	//		//	out <- in
	//		//}
	//
	//	}()
	//
	//	wg.Wait()
	//	return out
	//}

	//close(out)

	//return nil

	//if isFinished {
	//	return nil
	//}

}
