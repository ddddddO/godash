package main

import (
	"fmt"
	"time"
	"sync"
)

func main() {
	tasksCh := make(chan int)
	wg := &sync.WaitGroup{}

	// 複数タスク受け付けてキューにエンキューするgoroutine
	go recieveTasks(tasksCh, wg)
	wg.Add(1)

	// キューから受け付けたタスクをデキューして処理するgoroutine
	go processTasks(tasksCh, wg)
	wg.Add(1)

	// time.Sleep(13 * time.Second)
	wg.Wait()
	fmt.Println("done...")
}

func recieveTasks(tasks chan<- int, wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		tasks <- i
		time.Sleep(1 * time.Second)
	}

	close(tasks)
	wg.Done()
}

func processTasks(tasks <-chan int, wg *sync.WaitGroup) {
	for t := range tasks {
		fmt.Printf("task:%d\n", t)
	}

	wg.Done()
}