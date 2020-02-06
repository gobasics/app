package app

import "os"

type Option func(*server)

func WithAutoCert() Option {
	return func(s *server) {
		s.tlsProvider = &autoCert{}
	}
}

func WithServer(b Server) Option {
	return func(s *server) {
		s.backend = b
	}
}

func WithHost(ip string) Option {
	return func(s *server) {
		s.hostIP = ip
	}
}
func WithPort(port int) Option {
	return func(s *server) {
		s.port = port
	}
}

func WithStopChan(stopChan chan os.Signal) Option {
	return func(s *server) {
		s.stopChan = stopChan
	}
}
