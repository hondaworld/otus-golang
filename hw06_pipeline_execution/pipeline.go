package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func take(in In, done In, stage Stage) Out {
	takeStream := make(Bi)
	tempStream := make(Bi)

	go func() {
		defer close(takeStream)
		//defer wg.Done()

		select {
		case <-done:
			return
		case takeStream <- stage(tempStream):
		}
	}()

	go func() {
		defer close(tempStream)

		for v := range in {
			//wg.Add(1)
			tempStream <- v
		}
	}()

	return takeStream
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	outTemp := make(Bi)
	temp := make(Bi)
	wg := sync.WaitGroup{}
	//
	//wg.Add(1)

	//pipeline := concurrentPipeline(in, done, stages...)

	//go func() {
	//	defer close(outTemp)
	//
	//	select {
	//	case <-done:
	//		return
	//	case temp <- outTemp:
	//	}
	//}()

	go func() {
		wg.Wait()
		//out <- temp
	}()

	go func() {
		defer wg.Done()
		for _, stage := range stages {
			//wg.Add(1)
			//s := stage

			//go func() {

			//for {
			select {
			case <-done:
				break
			case temp <- <-outTemp:
			//case outTemp <- <-take(temp, done, stage):
			case outTemp <- <-stage(temp):
				//case temp <- outTemp:
			}
			//}
			//}()
		}
	}()

	go func() {
		defer close(temp)

		//for {
		for v := range in {

			select {
			case <-done:
				return
			case temp <- v:
				wg.Add(1)
				//case temp <- outTemp:
			}
		}

		//out <- temp
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
