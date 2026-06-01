package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	conn *net.UDPConn
}

const (
	PORT      = ":5555"
	CONN_TYPE = "udp"
)

func Serve() error {
	addr, err := net.ResolveUDPAddr(CONN_TYPE, PORT)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP(CONN_TYPE, addr)
	if err != nil {
		return err
	}
	s := &server{
		conn,
	}
	go s.listen()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	return nil
}

func (s *server) listen() {
	defer s.conn.Close()
	buff := make([]byte, 512)
	for {
		n, err := s.conn.Read(buff)
		if err != nil {
			log.Fatalf("fail")
		}
		if n > 0 {
			fmt.Printf("success here is the content %s", buff)
		}
	}
}
