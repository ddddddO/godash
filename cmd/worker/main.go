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

				ss := &dummySecretStore{}
				ds := &postgreSQL{}
				w := worker{
					ss: ss,
					ds: ds,
				}

				statusCode := 200
				queryResult, err := w.run(t.DataSourceType, t.Query)
				if err != nil {
					statusCode = 500
				}

				result := &model.Result{
					StatusCode:  statusCode,
					QueryResult: queryResult,
				}
				if err := json.NewEncoder(t.conn).Encode(result); err != nil {
					fmt.Println(err)
				}
			}()
		}
	}
}

type worker struct {
	ss secretStore // file or redash or postgres or embeded db or ...
	ds dataSource  // 増やせるだけ...
}

type secretStore interface {
	load(dataSourceType string) (interface{}, error)
}

type dataSource interface {
	parse(string) error
	connect(interface{}) error
	execute(string) (string, error)
	close() error
}

func (w *worker) run(typ, query string) (string, error) {
	secret, err := w.ss.load(typ)
	if err != nil {
		return "", err
	}

	if err := w.ds.parse(query); err != nil {
		return "", err
	}

	if err := w.ds.connect(secret); err != nil {
		return "", err
	}
	defer w.ds.close()

	result, err := w.ds.execute(query)
	if err != nil {
		return "", err
	}

	return result, nil
}

type dummySecretStore struct {
}

func (dummySecretStore) load(typ string) (interface{}, error) {
	fmt.Println("not yet impl")

	_ = typ
	return "dummy secret", nil
}

type postgreSQL struct {
	// having db connection
}

func (pq *postgreSQL) parse(query string) error {
	fmt.Println("not yet impl")
	fmt.Println(query)
	return nil
}

func (pq *postgreSQL) connect(raw interface{}) error {
	fmt.Println("not yet impl")

	secret := raw.(string)
	fmt.Println("connect to data source using secret", secret)
	return nil
}

func (pq *postgreSQL) execute(query string) (string, error) {
	fmt.Println("not yet impl")

	_ = query
	return "555", nil
}

func (pq *postgreSQL) close() error {
	fmt.Println("not yet impl")
	return nil
}
