package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks without errors with Eventually", func(t *testing.T) {
		tasksCount := 3
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for inc := 15; inc < 40; inc += 10 {
			taskSleep := time.Millisecond * time.Duration(inc)
			err := fmt.Errorf("error")
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				return err
			})
		}

		for inc := 10; inc < 50; inc += 10 {
			for i := 0; i < tasksCount; i++ {
				taskSleep := time.Millisecond * time.Duration(inc)
				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}
		}

		workersCount := 5
		_ = Run(tasks, workersCount, 1)

		require.Eventually(t, func() bool {
			return runTasksCount == 4
		}, time.Millisecond*10, time.Millisecond, "not all tasks were completed")

		runTasksCount = 0
		_ = Run(tasks, workersCount, 2)

		require.Eventually(t, func() bool {
			return runTasksCount == 6
		}, time.Millisecond*10, time.Millisecond, "not all tasks were completed")

		workersCount = 6
		runTasksCount = 0
		_ = Run(tasks, workersCount, 3)

		require.Eventually(t, func() bool {
			return runTasksCount == 11
		}, time.Millisecond*10, time.Millisecond, "not all tasks were completed")
	})

	t.Run("tasks without errors with no errors count", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
	})
}
