package main

import (
	"fmt"
	"net"
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
	if _, err := conn.Write([]byte("555")); err != nil {
		panic(err)
	}
}