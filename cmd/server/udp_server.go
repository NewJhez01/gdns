package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gdns/internal/infrastructure/parser"
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
		_, _, err := s.conn.ReadFromUDP(buff)
		if err != nil {
			log.Fatalf("fail")
		}
		headerBuff := [12]byte(buff[:12])
		h := parser.CreateHeaders()
		h.Parse(headerBuff)

	}
}
