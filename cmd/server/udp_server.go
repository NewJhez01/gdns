package server

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	"github.com/gdns/internal/blocklist"
	"github.com/gdns/internal/cache"
	"github.com/gdns/internal/dns"
)

type Server struct {
	Conn    *net.UDPConn
	rClient cache.RedisClient
	s       blocklist.SqliteClient
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewServer(redisClient cache.RedisClient, b blocklist.SqliteClient) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		rClient: redisClient,
		s:       b,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (srv *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", os.Getenv("UDP_PORT"))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	srv.Conn = conn
	srv.wg.Add(1)
	go srv.listen()
	return nil
}

func (srv *Server) Stop() {
	srv.cancel()
	if srv.Conn != nil {
		srv.Conn.Close()
	}
	srv.wg.Wait()
}

func (srv *Server) listen() {
	defer srv.wg.Done()
	buff := make([]byte, 512)
	for {
		select {
		case <-srv.ctx.Done():
			return
		default:
		}
		n, addr, err := srv.Conn.ReadFromUDP(buff)
		if err != nil {
			if srv.ctx.Err() != nil {
				return // graceful shutdown
			}
			log.Printf("failed to read from udp: %s", err)
			continue
		}
		resp, err := dns.Resolve(buff[:n], &srv.rClient, &srv.s)
		if err != nil {
			log.Printf("failed to resolve: %s", err)
			continue
		}
		if _, err := srv.Conn.WriteToUDP(resp, addr); err != nil {
			log.Printf("failed to write response: %s", err)
		}
	}
}
