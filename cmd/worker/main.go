package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ddddddO/godash/cmd/worker/datasource"
	"github.com/ddddddO/godash/cmd/worker/secretstore"
	"github.com/ddddddO/godash/model"

	"github.com/urfave/cli/v2"
)

type taskAndConn struct {
	*model.Task
	conn net.Conn
}

// TODO:
// データソース接続情報取得（どこから？DB or redis?）<- clientからもらってファイルに保存した
// データソース接続 <- postgresqlのみ対応した
// クエリ取得 <- clientから受信できるよう対応した
// クエリパース <- TODO: まだ
// クエリ投げる <- postgresのみ対応した
// クエリ結果の返却 <- TODO: 特定のクエリしか対応してない。どうすればいいんだろう
func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "worker run",
				Action:  action,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func action(c *cli.Context) error {
	log.Println("start worker")

	ctx, cancel := context.WithCancel(context.Background())
	tasks := make(chan *taskAndConn)

	// 複数タスク受け付けてキューにエンキューするgoroutine
	go recieveTasks(ctx, tasks)

	// キューから受け付けたタスクをデキューして処理するgoroutine
	go processTasks(ctx, tasks)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)

	<-sig
	log.Println("graceful shutdown...")
	cancel()
	time.Sleep(3 * time.Second)
	return nil
}

func recieveTasks(_ context.Context, tasks chan<- *taskAndConn) {
	defer func() {
		close(tasks)
	}()

	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Printf("cannot listen: %v\n", err)
		return
	}

	// 接続を待ち受け続ける
	for {
		// 1接続分
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("cannot accept: %v\n", err)
			continue
		}
		log.Println("connected")

		// 複数の接続を扱うためgoroutine
		go func() {
			log.Println("received task")

			receivedTask := &model.Task{}
			if err := json.NewDecoder(conn).Decode(receivedTask); err != nil {
				log.Println(err)
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

// clientから受信した接続情報を保存する先のファイルパス
// worker起動時に設定ファイルから接続情報格納先を取得して決める、でもいいかも
const filePath = "/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/godash/testdata/postgres_connection_info"

func processTasks(ctx context.Context, tasks <-chan *taskAndConn) {
	for {
		select {
		case <-ctx.Done():
			log.Println("process done")
			return
		case t := <-tasks:

			w := newWorker(
				secretstore.NewFile(filePath),
				datasource.NewPostgreSQL(),
			)

			switch t.Task.Kind {
			case model.KindSettings:
				go settings(t, w)
			case model.KindQuery:
				go query(t, w)
			default:
				panic("unknown task kind")
			}
		}
	}
}

func settings(t *taskAndConn, w *worker) {
	defer t.conn.Close()

	statusCode := 200
	err := w.runSettings(t.DataSourceType, t.Settings)
	if err != nil {
		log.Println(err)
		statusCode = 500
	}

	result := &model.Result{
		StatusCode: statusCode,
	}
	if err := json.NewEncoder(t.conn).Encode(result); err != nil {
		log.Println(err)
	}
}

func query(t *taskAndConn, w *worker) {
	defer t.conn.Close()

	statusCode := 200
	queryResult, err := w.runQuery(t.DataSourceType, t.Query)
	if err != nil {
		log.Println(err)
		statusCode = 500
	}

	result := &model.Result{
		StatusCode:  statusCode,
		QueryResult: queryResult,
	}
	if err := json.NewEncoder(t.conn).Encode(result); err != nil {
		log.Println(err)
	}
}
