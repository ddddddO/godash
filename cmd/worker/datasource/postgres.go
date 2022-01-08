package datasource

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v4"
)

type postgreSQL struct {
	conn        *pgx.Conn
	parsedQuery *parsedQuery
}

func NewPostgreSQL() *postgreSQL {
	return &postgreSQL{}
}

// TODO: ここをやっていく
// クエリ文字列の先頭の空白除去
// select/insert/update/deleteの文字列がプレフィックスにあれば、一旦パース成功とみなす
func (pg *postgreSQL) Parse(query string) error {
	pq := newParsedQuery(query)
	if err := pq.validate(); err != nil {
		return err
	}
	pq.decideQueryType()

	pg.parsedQuery = pq

	return nil
}

func (pg *postgreSQL) Connect(raw interface{}) error {
	url := raw.(string)

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return err
	}
	pg.conn = conn
	return nil
}

// TODO: ここをやっていく
// 難しそう。ParseメソッドでExecuteメソッドが使いやすいようなstructを用意した方がいいかも
func (pg *postgreSQL) Execute() (string, error) {
	rows, err := pg.conn.Query(context.TODO(), pg.parsedQuery.value)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	ret := ""
	onceGetColumns := sync.Once{}
	for rows.Next() {
		// 最初だけ選択されたカラム名を取得
		onceGetColumns.Do(
			func() {
				fds := rows.FieldDescriptions()
				header := ""
				for _, fd := range fds {
					header += string(fd.Name) + " "
				}

				ret += header + fmt.Sprintln()
			},
		)

		// これをつかえば良さそう
		// https://pkg.go.dev/github.com/jackc/pgx#Rows.Values
		values, err := rows.Values()
		if err != nil {
			fmt.Println(err)
			return "", err
		}

		// タイプアサーションの数を増やせばよさそう
		for _, v := range values {
			switch v.(type) {
			case string:
				ret += v.(string) + " "
			case int:
				ret += string(strconv.Itoa(v.(int)))
			}
		}
		ret += fmt.Sprintln()
	}

	fmt.Println(ret)

	return ret, nil
}

func (pg *postgreSQL) Close() error {
	return pg.conn.Close(context.TODO())
}
