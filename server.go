package app

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
)

type Server interface {
	GracefulStop()
	Serve(net.Listener) error
}

type TLSProvider interface {
	TLSConfig() (*tls.Config, error)
}

type App struct {
	hostIP   string
	listener net.Listener
	port     int

	backend     Server
	stopChan    chan os.Signal
	tlsProvider TLSProvider
}

func (a *App) listen() error {
	var err error
	hostPort := fmt.Sprintf("%s:%d", a.hostIP, a.port)

	if a.tlsProvider == nil {
		a.listener, err = net.Listen("tcp", hostPort)
	} else {
		var tlsConfig *tls.Config

		if tlsConfig, err = a.tlsProvider.TLSConfig(); err != nil {
			return err
		}

		a.listener, err = tls.Listen("tcp", hostPort, tlsConfig)
	}

	if err != nil {
		return fmt.Errorf("could not listen on %s; %w", hostPort, err)
	}

	return nil
}

func (a *App) serve() error {
	errChan := make(chan error)

	go func() {
		errChan <- a.backend.Serve(a.listener)
	}()

	select {
	case err := <-errChan:
		return err
	case _ = <-a.stopChan:
		a.backend.GracefulStop()
		return nil
	}
}

func (a *App) Start() error {
	if err := a.listen(); err != nil {
		return err
	}

	signal.Notify(a.stopChan, os.Interrupt, os.Kill)

	return a.serve()
}

func New(options ...Option) *App {
	a := App{stopChan: make(chan os.Signal, 1)}

	for _, f := range options {
		f(&a)
	}

	return &a
}
