package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// クエリ結果の返却

	fmt.Println("start worker")

	ctx, cancel := context.WithCancel(context.Background())
	run(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)

	<-sig
	fmt.Println("graceful shutdown...")
	cancel()
	time.Sleep(3 * time.Second)
}

func run(ctx context.Context) {
	tasks := make(chan *taskAndConn)

	// 複数タスク受け付けてキューにエンキューするgoroutine
	go recieveTasks(ctx, tasks)

	// キューから受け付けたタスクをデキューして処理するgoroutine
	go processTasks(ctx, tasks)
}


func recieveTasks(_ context.Context, tasks chan<- *taskAndConn) {
	defer func() {
		close(tasks)
	}()

	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("cannot listen", err)
		return
	}

	// 接続を待ち受け続ける
	for {
		// 1接続分
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("cannot accept", err)
			continue
		}
		fmt.Println("connected")

		// 複数の接続を扱うためgoroutine
		go func() {
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
}

func processTasks(ctx context.Context, tasks <-chan *taskAndConn) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("process done")
			return
		case t := <-tasks:
			go func() {
				defer t.conn.Close()

				fmt.Printf("Task\ndata source type: %s\nquery: %s\n", t.DataSourceType, t.Query)

				result := &model.Result{
					StatusCode:  200,
					QueryResult: "xxxxx",
				}
				if err := json.NewEncoder(t.conn).Encode(result); err != nil {
					fmt.Println(err)
				}
			}()
		}
	}
}
