package app

type Option func(*Server)

func WithAutoCert() Option {
	return func(s *Server) {
		s.tlsProvider = &autoCert{}
	}
}

func WithBackend(b Backend) Option {
	return func(s *Server) {
		s.backend = b
	}
}

func WithHostIP(ip string) Option {
	return func(s *Server) {
		s.hostIP = ip
	}
}
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}
