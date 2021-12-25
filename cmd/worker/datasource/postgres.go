package datasource

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

type postgreSQL struct {
	// having db connection
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
	const url = "postgres://postgres:passw0rd@localhost:15432/dvdrental"

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	fmt.Println("not yet impl")

	secret := raw.(string)
	fmt.Println("connect to data source using secret", secret)
	return nil
}

func (pg *postgreSQL) Execute(query string) (string, error) {
	fmt.Println("not yet impl")

	_ = query
	return "555", nil
}

func (pg *postgreSQL) Close() error {
	fmt.Println("not yet impl")
	return nil
}
