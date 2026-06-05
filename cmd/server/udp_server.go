package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gdns/internal/dns/parser"
)

type server struct {
	conn *net.UDPConn
}

const (
	PORT      = ":5555"
	CONN_TYPE = "udp"
)

func Serve() {
	addr, err := net.ResolveUDPAddr(CONN_TYPE, PORT)
	if err != nil {
		log.Fatalf("failed to resolve udp addr prev: %s", err)
	}
	conn, err := net.ListenUDP(CONN_TYPE, addr)
	if err != nil {
		log.Fatalf("failed to listen to udp addr prev: %s", err)
	}
	s := &server{
		conn,
	}
	go s.listen()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func (s *server) listen() {
	defer s.conn.Close()
	buff := make([]byte, 512)
	for {
		_, _, err := s.conn.ReadFromUDP(buff)
		if err != nil {
			log.Fatalf("fail")
		}
		dns := parser.CreateNewDnsStruct()
		dns.Parse([512]byte(buff))
	}
}
