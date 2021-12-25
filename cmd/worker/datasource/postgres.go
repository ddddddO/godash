package datasource

import (
	"context"
	"fmt"
	"os"

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
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	pg.conn = conn

	fmt.Println("connect to data source using secret")
	return nil
}

func (pg *postgreSQL) Execute(query string) (string, error) {
	query = `select tablename, tableowner from pg_catalog.pg_tables where schemaname = 'public'`
	rows, err := pg.conn.Query(context.TODO(), query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var owner string
		rows.Scan(&name, &owner)
		fmt.Printf("%s owned by %s\n", name, owner)
	}

	return "555", nil
}

func (pg *postgreSQL) Close() error {
	fmt.Println("not yet impl")
	return nil
}
