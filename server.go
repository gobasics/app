package app

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

type Provider interface {
	GracefulStop()
	Serve(net.Listener) error
}

type Server struct {
	Config   Config
	Provider Provider

	listener  net.Listener
	stop      chan os.Signal
	tlsConfig *tls.Config
}

func (s *Server) parseDirCache() error {
	if len(s.Config.DirCache) == 0 {
		return errors.New("autocert certificates cache dir is not set")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1<<32 - 1)

	_, err := os.Create(fmt.Sprintf("%s%d", s.Config.DirCache, r))

	return err

}

func (s *Server) parseHosts() error {
	if len(s.Config.HostNames) == 0 {
		return errors.New("hosts whitelist is not set")
	}

	for i := range s.Config.HostNames {
		s.Config.HostNames[i] = strings.TrimSpace(s.Config.HostNames[i])
	}

	return nil
}

func (s *Server) initLetsencrypt() error {
	if err := s.parseDirCache(); err != nil {
		return err
	}

	if err := s.parseHosts(); err != nil {
		return err
	}

	m := &autocert.Manager{
		Cache:      autocert.DirCache(s.Config.DirCache),
		HostPolicy: autocert.HostWhitelist(s.Config.HostNames...),
		Prompt:     autocert.AcceptTOS,
	}

	s.tlsConfig = m.TLSConfig()

	return nil
}

func (s *Server) listen() error {
	var err error

	hostPort := fmt.Sprintf(":%d", s.Config.Port)
	if s.tlsConfig == nil {
		s.listener, err = net.Listen("tcp", hostPort)
	} else {
		s.listener, err = tls.Listen("tcp", hostPort, s.tlsConfig)
	}

	if err != nil {
		return fmt.Errorf("could not listen on %s; %w", hostPort, err)
	}

	return nil
}

func (s *Server) serve() error {
	return s.Provider.Serve(s.listener)
}

func (s *Server) Done() chan os.Signal {
	go func() {
		s.Provider.GracefulStop()
	}()

	return s.stop
}

func (s *Server) ListenAndServe() error {
	if s.Config.Letsencrypt {
		if err := s.initLetsencrypt(); err != nil {
			return err
		}
	}

	if err := s.listen(); err != nil {
		return err
	}

	return s.serve()
}
