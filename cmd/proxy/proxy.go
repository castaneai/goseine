package main

import (
	"log"
	"net"
	"os"

	"github.com/castaneai/goseine/proxy"
)

func main() {
	loginServerAddr := os.Getenv("GOSEINE_LOGIN_SERVER")
	if loginServerAddr == "" {
		log.Fatalf("env: GOSEINE_LOGIN_SERVER is undefined")
	}
	pool := proxy.NewPool()
	login := proxy.NewLoginProxy(loginServerAddr, pool)

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:20000")
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	login.Start(lis)
}
