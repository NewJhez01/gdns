package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gdns/internal/cache"
	"github.com/gdns/internal/dns"
)

type server struct {
	conn    *net.UDPConn
	rClient cache.RedisClient
}

const (
	PORT      = ":5555"
	CONN_TYPE = "udp"
)

func Serve(redisClient cache.RedisClient) {
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
		redisClient,
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
		n, addr, err := s.conn.ReadFromUDP(buff)
		if err != nil {
			log.Print("failed to read from udp continue")
			continue
		}
		resp, err := dns.Resolve(buff[:n], &s.rClient)
		if err != nil {
			log.Print("failed to resolve the dns request continue")
			continue
		}
		s.conn.WriteToUDP([]byte(resp), addr)
	}
}
