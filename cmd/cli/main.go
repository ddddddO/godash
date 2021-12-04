package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ddddddO/godash/model"
)

const (
	protocol   = "tcp"
	targetHost = "localhost"
	targetPort = 9999
)

func main() {
	conn, err := net.Dial(
		protocol,
		fmt.Sprintf("%s:%d", targetHost, targetPort),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("send task")

	pgTask := model.Task{
		DataSourceType: "postgres",
		Query:          "select * from test",
	}

	if err := json.NewEncoder(conn).Encode(pgTask); err != nil {
		panic(err)
	}
}
