package main

import (
	"os"
	"net"
	"fmt"
	"github.com/castaneai/goseine"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10000")
	if err != nil {
		fmt.Println(err)
		return
	}
	actual, err := net.ResolveTCPAddr("tcp", os.Getenv("GOSEINE_ACTUAL_LOGIN_SERVER"))
	if err != nil {
		fmt.Println(err)
		return
	}
	goseine.Start(addr, actual)
}
