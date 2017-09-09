package main

import (
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
	actual, err := net.ResolveTCPAddr("tcp", "157.7.172.197:10000")
	if err != nil {
		fmt.Println(err)
		return
	}
	goseine.Start(addr, actual)
}
