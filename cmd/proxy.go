package main

import (
	"log"
	"net"

	"github.com/castaneai/goseine"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Fatal(err)
	}
	/*
		actual, err := net.ResolveTCPAddr("tcp", os.Getenv("GOSEINE_ACTUAL_LOGIN_SERVER"))
		if err != nil {
			log.Fatal(err)
		}
	*/

	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := goseine.NewServer()
	log.Fatal(s.Serve(lis))
}
