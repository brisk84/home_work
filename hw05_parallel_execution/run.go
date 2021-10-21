package hw05parallelexecution

import (
	"errors"
    "fmt"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

    //for _, task := range tasks {
    //var workersCount = 5
    var tasksCount = len(tasks)
    var errCount = 0
    var curTask = 0
    fmt.Printf("tasks: %d, n: %d, m: %d\n", tasksCount, n, m)

    c := make(chan error, n)
    for {
        task := tasks[curTask]
        go func(task Task) {
            c <- task()
        }(task)
        err := <-c
        if err != nil {
            errCount++
        }
        curTask++
        if (curTask >= tasksCount) || (errCount >= m) {
            break
        }
    }
    fmt.Printf("Tasks completed: %d\n", curTask)
    if errCount >= m {
        return ErrErrorsLimitExceeded
    }
/*
    for i := 0; i < workersCount; i++ {
        c := make(chan error, n)

        task := tasks[curTask]
        go func(Task) {
            err := tasks[i]()
            c <- err
        }(task)
        //err := task()
        fmt.Println(<-c)
    }
*/
    //for i := 1; i < n; i++ {
        //go func(i int) {
        //    fmt.Println("Test ", i)
        //}(i)
    //}
    //fmt.Println(tasks, n, m)
    fmt.Println()

	return nil
}
