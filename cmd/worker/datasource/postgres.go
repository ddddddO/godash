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
	rows, err := pg.conn.Query(context.TODO(), query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var firstName string
		var lastName string
		rows.Scan(&firstName, &lastName)
		fmt.Printf("%s %s\n", firstName, lastName)
	}

	return "555", nil
}

func (pg *postgreSQL) Close() error {
	return pg.conn.Close(context.TODO())
}
