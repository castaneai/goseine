package main

import (
	"log"
	"net"
	"os"

	"github.com/castaneai/goseine"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:20000")
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	loginServerAddr := os.Getenv("GOSEINE_LOGIN_SERVER")
	if loginServerAddr == "" {
		log.Fatalf("env: GOSEINE_LOGIN_SERVER is undefined")
	}

	dialer := goseine.NewProxyDialer(os.Getenv("GOSEINE_LOGIN_SERVER"))
	p := goseine.NewProxy(dialer)
	log.Fatal(p.Start(lis))
}
