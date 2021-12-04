package main

import (
	"encoding/json"
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
		Query:          "select * from test",
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// query結果を受け取るよう
	go func() {
		result := make([]byte, 9)
		if _, err := conn.Read(result); err != nil {
			panic(err)
		}

		fmt.Println(string(result))
		wg.Done()
	}()

	// taskをworkerプロセスへ
	if err := json.NewEncoder(conn).Encode(pgTask); err != nil {
		panic(err)
	}

	wg.Wait()
}
