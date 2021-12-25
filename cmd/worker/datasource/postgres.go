package datasource

import (
	"fmt"
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
