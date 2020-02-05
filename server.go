package app

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Backend interface {
	GracefulStop()
	Serve(net.Listener) error
}

type TLSProvider interface {
	TLSConfig() (*tls.Config, error)
}

type server struct {
	hostIP   string
	listener net.Listener
	port     int

	backend     Backend
	stopChan    chan os.Signal
	tlsProvider TLSProvider
}

func (s *server) listen() error {
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

func (s *server) serve() error {
	errChan := make(chan error)

	go func() {
		errChan <- s.backend.Serve(s.listener)
	}()

	select {
	case err := <-errChan:
		return err
	case <-s.stopChan:
		s.backend.GracefulStop()
		return nil
	}
}

func Run(options ...Option) error {
	s := server{stopChan: make(chan os.Signal)}

	signal.Notify(s.stopChan, syscall.SIGINT, syscall.SIGTERM)

	for _, f := range options {
		f(&s)
	}

	if err := s.listen(); err != nil {
		return err
	}

	return s.serve()
}
