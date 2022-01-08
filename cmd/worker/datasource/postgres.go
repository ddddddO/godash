package datasource

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type queryType uint

const (
	undefined queryType = iota
	selectType
	insertType
	updateType
	deleteType
)

type parsedQuery struct {
	qType queryType
	query string
}

type postgreSQL struct {
	conn        *pgx.Conn
	parsedQuery *parsedQuery
}

func NewPostgreSQL() *postgreSQL {
	return &postgreSQL{}
}

var (
	errUndefinedType = errors.New("sql query is undefined dml")
)

// TODO: ここをやっていく
// クエリ文字列の先頭の空白除去
// select/insert/update/deleteの文字列がプレフィックスにあれば、一旦パース成功とみなす
func (pg *postgreSQL) Parse(query string) error {
	q := strings.TrimSpace(query)
	if err := validate(q); err != nil {
		return err
	}

	pg.parsedQuery = &parsedQuery{
		qType: decideQueryType(q),
		query: q,
	}

	return nil
}

func validate(query string) error {
	switch {
	case strings.HasPrefix(query, "select"):
		return nil
	case strings.HasPrefix(query, "insert"):
		return nil
	case strings.HasPrefix(query, "update"):
		return nil
	case strings.HasPrefix(query, "delete"):
		return nil
	}
	return errUndefinedType
}

func decideQueryType(query string) queryType {
	switch {
	case strings.HasPrefix(query, "select"):
		return selectType
	case strings.HasPrefix(query, "insert"):
		return insertType
	case strings.HasPrefix(query, "update"):
		return updateType
	case strings.HasPrefix(query, "delete"):
		return deleteType
	}
	return undefined
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
	rows, err := pg.conn.Query(context.TODO(), pg.parsedQuery.query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	ret := ""
	for rows.Next() {
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
				// FIXME: using strconv
				ret += string(v.(int))
			}
		}

		fmt.Println(ret)
	}

	return ret, nil
}

func (pg *postgreSQL) Close() error {
	return pg.conn.Close(context.TODO())
}
