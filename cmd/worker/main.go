package main

import (
	"fmt"
	"net"
	"sync"
	"strconv"
	"io"
)

func main() {
	fmt.Println("start worker")

	tasksCh := make(chan int)
	wg := &sync.WaitGroup{}

	// 複数タスク受け付けてキューにエンキューするgoroutine
	go recieveTasks(tasksCh, wg)
	wg.Add(1)

	// キューから受け付けたタスクをデキューして処理するgoroutine
	go processTasks(tasksCh, wg)
	wg.Add(1)

	wg.Wait()
	fmt.Println("done...")
}

func recieveTasks(tasks chan<- int, wg *sync.WaitGroup) {
	// NOTE: unixドメインソケットがwslでは使えないよう
	// ln, err := net.Listen("unix", "./sock")
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("cannot listen", err)
	}

	// 接続を待ち受け続ける
	for {
		// 1接続分
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("cannot accept", err)
		}
		fmt.Println("connected")

		// 複数の接続を扱うためgoroutine
		go func() {
			defer conn.Close()

			buf := make([]byte, 3) // NOTE: 送られてくるデータのサイズに合わせないと余計なパディング含まれる
			_, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("connection closed...")
					return
				}

				fmt.Println("cannot read", err)
				return
			}

			fmt.Println("received task")

			t, err := strconv.Atoi(string(buf))
			if err != nil {
				fmt.Println("can not cast", err)
				return
			}

			tasks <- t
		}()
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