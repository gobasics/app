package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"

	"gobasics.dev/env"
)

type Provider interface {
	GracefulStop()
	Serve(net.Listener) error
}

type TLSProvider interface {
	TLSConfig() (*tls.Config, error)
}

type Server struct {
	hostIP   string
	listener net.Listener
	port     int

	provider    Provider
	stopChan    chan os.Signal
	tlsProvider TLSProvider
}

func (s *Server) listen() error {
	var err error
	hostPort := fmt.Sprintf("%s:%d", s.hostIP, s.port)

	if s.tlsProvider == nil {
		s.listener, err = net.Listen("tcp", hostPort)
	} else {
		var tlsConfig *tls.Config

		if tlsConfig, err = s.tlsProvider.TLSConfig(); err != nil {
			return err
		}

		s.listener, err = tls.Listen("tcp", hostPort, tlsConfig)
	}

	if err != nil {
		return fmt.Errorf("could not listen on %s; %w", hostPort, err)
	}

	return nil
}

func (s *Server) serve() error {
	signal.Notify(s.stopChan, os.Interrupt, os.Kill)

	go func() {
		<-s.stopChan
		s.provider.GracefulStop()
	}()

	return s.provider.Serve(s.listener)
}

func (s *Server) Start() error {
	if err := s.listen(); err != nil {
		return err
	}

	return s.serve()
}

func DefaultOptions() []Option {
	port, _ := env.Get("PORT").Int()
	dirCache := env.Get("DIR_CACHE").Str()
	hostnames := env.Get("HOSTNAMES").StrSlice(",")
	return []Option{WithPort(port), WithAutoCert(dirCache, hostnames...)}
}

func New(options ...Option) *Server {
	s := Server{stopChan: make(chan os.Signal, 1)}

	for _, o := range options {
		o(&s)
	}

	return &s
}
