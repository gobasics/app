package app

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
)

type Backend interface {
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

	backend     Backend
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
	return s.backend.Serve(s.listener)
}

func (s *Server) Done() chan os.Signal {
	stop := make(chan os.Signal, 1)

	go func() {
		<-stop
		s.backend.GracefulStop()
	}()

	return stop
}

func (s *Server) ListenAndServe() error {
	if err := s.listen(); err != nil {
		return err
	}

	return s.serve()
}

func New(options ...Option) *Server {
	s := Server{}

	for _, f := range options {
		f(&s)
	}

	return &s
}
