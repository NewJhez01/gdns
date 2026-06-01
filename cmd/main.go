package main

import (
	"fmt"
	"log"
	"net"
)

const (
	PORT = "5555"
	TYPE = "tcp"
)

func main() {
	l, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("failed to listen to port prev: %s", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Failed to accept request prev: %s", err)
		}
		go func() {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				log.Fatalf("Failed to read from request prev: err")
			}
			if n != 0 {
				fmt.Printf("success here is the content: %s", string(buff))
			}
			defer conn.Close()
		}()
	}
}
