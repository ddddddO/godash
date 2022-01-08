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
	selectType queryType = iota
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
	if !strings.HasPrefix(query, "select") {
		return errUndefinedType
	}

	pg.parsedQuery = &parsedQuery{
		qType: selectType,
		query: query,
	}

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
	rows, err := pg.conn.Query(context.TODO(), pg.parsedQuery.query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	ret := ""
	for rows.Next() {
		var firstName string
		var lastName string
		rows.Scan(&firstName, &lastName)
		ret = fmt.Sprintf("%s %s\n", firstName, lastName)
	}

	return ret, nil
}

func (pg *postgreSQL) Close() error {
	return pg.conn.Close(context.TODO())
}
