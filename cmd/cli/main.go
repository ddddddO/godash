package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/ddddddO/godash/model"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

const (
	protocol   = "tcp"
	targetHost = "localhost"
	targetPort = 9999
)

// TODO:
// データソース接続情報（コマンドライン引数 or 設定ファイル or ...） <- コマンドライン引数から受け付けるよう対応した。
// データソース接続確認 <- worker(server)側で接続確認した結果がこちらに返ってくるようになってる
// データソース接続情報保存（どこに？DB or redis?） <- worker(server)側で接続確認してOKであれば、worker(server)側で保存してる
// クエリ取得（コマンドライン引数 or 標準入力 or ...） <- コマンドライン引数から受け付けるよう対応した。
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

		fmt.Printf("status: %d\n", result.StatusCode)
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

		fmt.Printf("status: %d\nquery result:\n", result.StatusCode)
		if err := showQueryResult(result.QueryResult); err != nil {
			fmt.Println(err)
			return
		}
	}()

	// taskをworkerプロセスへ
	if err := json.NewEncoder(conn).Encode(task); err != nil {
		panic(err)
	}

	wg.Wait()
	return nil
}

// FIXME: ここはworker側がどんな文字列を返すかによって処理が決まる。
// FIXME: あと表示がバグってる
func showQueryResult(raw string) error {
	rows := strings.Split(raw, "\n")
	header := strings.Split(rows[0], " ")
	valuesRows := rows[1:]
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	for _, row := range valuesRows {
		values := strings.Split(row, " ")
		if values[0] == "" { // FIXME: もっとちゃんとした方が良さそう
			continue
		}
		table.Append(values)
	}
	table.Render()
	return nil
}
