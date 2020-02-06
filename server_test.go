package app

import (
	"errors"
	"fmt"
	"net"
	"os"
	"testing"
)

type fakeBackend struct {
	serveErr              error
	stopChan              chan struct{}
	stopCount, serveCount int
}

func (fb *fakeBackend) GracefulStop() {
	fb.stopCount += 1
	fb.stopChan <- struct{}{}
}

func (fb *fakeBackend) Serve(l net.Listener) error {
	fb.serveCount += 1
	if fb.serveErr != nil {
		return fb.serveErr
	}
	<-fb.stopChan
	return nil
}

func TestGracefulStop(t *testing.T) {
	for _, test := range []struct {
		backend   *fakeBackend
		signal    os.Signal
		stopCount int
	}{
		{&fakeBackend{}, os.Interrupt, 1},
		{&fakeBackend{}, os.Kill, 1},
	} {
		t.Run(fmt.Sprintf("%+v", test), func(t *testing.T) {
			stopChan := make(chan os.Signal)
			want := test.stopCount
			go func() {
				_ = Run(WithBackend(test.backend), WithStopChan(stopChan))
			}()
			stopChan <- test.signal
			got := test.backend.stopCount
			if got != want {
				t.Errorf("wanted %v, got %v", want, got)
			}
		})
	}
}

func TestServeError(t *testing.T) {
	for _, test := range []struct {
		backend *fakeBackend
	}{
		{&fakeBackend{serveErr: errors.New("foo")}},
		{&fakeBackend{serveErr: errors.New("bar")}},
	} {
		t.Run(fmt.Sprintf("%+v", test), func(t *testing.T) {
			want := test.backend.serveErr
			got := Run(WithBackend(test.backend))
			if got != want {
				t.Errorf("wanted %v, got %v", want, got)
			}
		})
	}
}
