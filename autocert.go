package app

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

type autoCert struct {
	DirCache  string
	HostNames []string
}

func (ac autoCert) parseDirCache() error {
	if len(ac.DirCache) == 0 {
		return errors.New("autocert certificates cache dir is not set")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1<<32 - 1)

	_, err := os.Create(fmt.Sprintf("%s%d", ac.DirCache, r))

	return err

}

func (ac autoCert) parseHosts() error {
	if len(ac.HostNames) == 0 {
		return errors.New("hosts whitelist is not set")
	}

	return nil
}

func (ac *autoCert) init() error {
	if err := ac.parseDirCache(); err != nil {
		return err
	}

	if err := ac.parseHosts(); err != nil {
		return err
	}

	return nil
}

func (ac *autoCert) TLSConfig() (*tls.Config, error) {
	if err := ac.init(); err != nil {
		return nil, err
	}
	m := &autocert.Manager{
		Cache:      autocert.DirCache(ac.DirCache),
		HostPolicy: autocert.HostWhitelist(ac.HostNames...),
		Prompt:     autocert.AcceptTOS,
	}

	return m.TLSConfig(), nil
}
