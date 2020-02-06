package app

import "os"

type Option func(*App)

func WithAutoCert() Option {
	return func(a *App) {
		a.tlsProvider = &autoCert{}
	}
}

func WithServer(b Server) Option {
	return func(a *App) {
		a.backend = b
	}
}

func WithHost(ip string) Option {
	return func(a *App) {
		a.hostIP = ip
	}
}
func WithPort(port int) Option {
	return func(a *App) {
		a.port = port
	}
}

func WithStopChan(stopChan chan os.Signal) Option {
	return func(a *App) {
		a.stopChan = stopChan
	}
}
