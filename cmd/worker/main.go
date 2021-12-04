package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/ddddddO/godash/model"
)

type taskAndConn struct {
	*model.Task
	conn net.Conn
}

func main() {
	// TODO:
	// データソース接続情報取得（どこから？DB or redis?）
	// データソース接続
	// クエリ取得
	// クエリパース
	// クエリ投げる
	// クエリ結果の返却をどうしよう。cliの方に伝えないといけない

	fmt.Println("start worker")

	tasks := make(chan *taskAndConn)
	wg := &sync.WaitGroup{}

	// 複数タスク受け付けてキューにエンキューするgoroutine
	wg.Add(1)
	go recieveTasks(tasks, wg)

	// キューから受け付けたタスクをデキューして処理するgoroutine
	wg.Add(1)
	go processTasks(tasks, wg)

	wg.Wait()
	fmt.Println("done...")
}

func recieveTasks(tasks chan<- *taskAndConn, wg *sync.WaitGroup) {
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

			fmt.Println("received task")

			receivedTask := &model.Task{}
			if err := json.NewDecoder(conn).Decode(receivedTask); err != nil {
				fmt.Println(err)
				return
			}

			task := &taskAndConn{
				Task: receivedTask,
				conn: conn,
			}

			tasks <- task
		}()
	}

	close(tasks)
	wg.Done()
}

func processTasks(tasks <-chan *taskAndConn, wg *sync.WaitGroup) {
	for t := range tasks {
		fmt.Printf("Task\ndata source type: %s\nquery: %s\n", t.DataSourceType, t.Query)
	}

	wg.Done()
}
