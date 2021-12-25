package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"sync"

	"github.com/ddddddO/godash/model"
)

const (
	protocol   = "tcp"
	targetHost = "localhost"
	targetPort = 9999
)

func main() {
	// TODO:
	// データソース接続情報（コマンドライン引数 or 設定ファイル or ...）
	// データソース接続確認
	// データソース接続情報保存（どこに？DB or redis?）
	// クエリ取得（コマンドライン引数 or 標準入力 or ...）

	var (
		query string
	)
	flag.StringVar(&query, "q", "", "QUERY")
	flag.Parse()

	conn, err := net.Dial(
		protocol,
		fmt.Sprintf("%s:%d", targetHost, targetPort),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("send task")

	pgTask := model.Task{
		DataSourceType: "postgres",
		Query:          query,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// query結果を受け取るよう
	go func() {
		defer wg.Done()

		result := &model.Result{}
		if err := json.NewDecoder(conn).Decode(result); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Result\nstatus: %d\nquery result: %s\n", result.StatusCode, result.QueryResult)
	}()

	// taskをworkerプロセスへ
	if err := json.NewEncoder(conn).Encode(pgTask); err != nil {
		panic(err)
	}

	wg.Wait()
}
