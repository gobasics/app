package server

import "os"

type Option func(*Server)

func WithAutoCert(dirCache string, hostWhitelist ...string) Option {
	return func(s *Server) {
		s.tlsProvider = &autoCert{
			DirCache:      dirCache,
			HostWhitelist: hostWhitelist,
		}
	}
}

func WithServer(b Provider) Option {
	return func(s *Server) {
		s.provider = b
	}
}

func WithHost(ip string) Option {
	return func(s *Server) {
		s.hostIP = ip
	}
}
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithStopChan(stopChan chan os.Signal) Option {
	return func(s *Server) {
		s.stopChan = stopChan
	}
}
