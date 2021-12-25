package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/ddddddO/godash/model"
	"github.com/urfave/cli/v2"
)

const (
	protocol   = "tcp"
	targetHost = "localhost"
	targetPort = 9999
)

// TODO:
// データソース接続情報（コマンドライン引数 or 設定ファイル or ...）
// データソース接続確認
// データソース接続情報保存（どこに？DB or redis?）
// クエリ取得（コマンドライン引数 or 標準入力 or ...）
func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "settings",
				Aliases: []string{"s"},
				Usage:   "data source settings send to worker for saving.",
				Action:  actionSettings,
			},
			{
				Name:    "query",
				Aliases: []string{"q"},
				Usage:   "query send to worker for executing.",
				Action:  actionQuery,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func actionSettings(c *cli.Context) error {
	conn, err := net.Dial(
		protocol,
		fmt.Sprintf("%s:%d", targetHost, targetPort),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	settings := c.Args().First()

	task := model.Task{
		Kind:           model.KindSettings,
		DataSourceType: "postgres",
		Settings:       settings,
	}

	fmt.Println("send task")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// 結果を受け取るよう
	go func() {
		defer wg.Done()

		result := &model.Result{}
		if err := json.NewDecoder(conn).Decode(result); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Result\nstatus: %d\n", result.StatusCode)
	}()

	// taskをworkerプロセスへ
	if err := json.NewEncoder(conn).Encode(task); err != nil {
		panic(err)
	}

	wg.Wait()
	return nil
}

func actionQuery(c *cli.Context) error {
	conn, err := net.Dial(
		protocol,
		fmt.Sprintf("%s:%d", targetHost, targetPort),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	query := c.Args().First()

	task := model.Task{
		Kind:           model.KindQuery,
		DataSourceType: "postgres",
		Query:          query,
	}

	fmt.Println("send task")

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
	if err := json.NewEncoder(conn).Encode(task); err != nil {
		panic(err)
	}

	wg.Wait()
	return nil
}
