package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(task Task, c chan error, q chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case c <- task():
			return
		case <-q:
			return
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var errCount int
	var curTask int
	var activeWorkers int
	var tasksCompleted int
	var mu sync.Mutex
	var wg sync.WaitGroup
	tasksCount := len(tasks)

	c := make(chan error, n)
	q := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := <-c
			mu.Lock()
			activeWorkers--
			tasksCompleted++
			if err != nil {
				errCount++
			}
			if (errCount >= m) || (tasksCompleted >= tasksCount) {
				mu.Unlock()
				close(q)
				break
			}
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			mu.Lock()
			if activeWorkers >= n {
				mu.Unlock()
				continue
			}
			mu.Unlock()
			task := tasks[curTask]
			wg.Add(1)
			go worker(task, c, q, &wg)
			mu.Lock()
			curTask++
			activeWorkers++
			if (curTask >= tasksCount) || (errCount >= m) {
				mu.Unlock()
				break
			}
			mu.Unlock()
		}
	}()
	<-q
	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
