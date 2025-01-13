package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task, len(tasks))
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	errorsCount := 0

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			isFinish := false

			for job := range jobs {
				mu.Lock()
				if m > 0 && errorsCount >= m {
					isFinish = true
				}
				mu.Unlock()

				if isFinish {
					break
				}

				err := job()

				if err != nil {
					mu.Lock()
					errorsCount++
					mu.Unlock()
				}
			}
		}()
	}

	for _, task := range tasks {
		jobs <- task
	}

	close(jobs)

	wg.Wait()

	if m > 0 && errorsCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
