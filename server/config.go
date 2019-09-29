package server

type Config struct {
	Letsencrypt bool
	DirCache    string
	HostNames   []string
	HostIP      string
	Port        int
}