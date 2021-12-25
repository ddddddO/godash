package datasource

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type postgreSQL struct {
	conn *pgx.Conn
}

func NewPostgreSQL() *postgreSQL {
	return &postgreSQL{}
}

func (pg *postgreSQL) Parse(query string) error {
	fmt.Println("not yet impl")
	fmt.Println(query)
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

func (pg *postgreSQL) Execute(query string) (string, error) {
	rows, err := pg.conn.Query(context.TODO(), query)
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
