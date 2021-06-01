package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"
)

func functionOptionsDemo() {
	//server, err := NewServer(":", 8080, Timeout(5*time.Second), Protocol("tcp"), MaxConns(100))
	//if err != nil{
	//	log.Fatal("server init failed!")
	//}
	handle := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Recieved Request %s from %s\n", r.URL.Path, r.RemoteAddr)
		fmt.Fprintf(w, "Hello, World! "+r.URL.Path)
	}
	http.HandleFunc("/v1/hello", handle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

type Option func(*Server)

func NewServer(addr string, port int, options ...Option) (*Server, error) {
	server := &Server{
		Addr: addr,
		Port: port,
	}
	for _, option := range options {
		option(server)
	}
	return server, nil
}

func Protocol(protocol string) Option {
	return func(server *Server) {
		server.Protocol = protocol
	}
}

func Timeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.Timeout = timeout
	}
}

func MaxConns(maxConns int) Option {
	return func(server *Server) {
		server.MaxConns = maxConns
	}
}

func Tls(tls *tls.Config) Option {
	return func(server *Server) {
		server.TLS = tls
	}
}
