package wrapper

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"igops.me/icap/server"
	"net/http"
	"os"
	"testing"
	"time"
)

type DummyICAPHandler struct {
}

func (h *DummyICAPHandler) OnREQMOD(request []byte) (response []byte, err error) {
	return request, nil
}

func (h *DummyICAPHandler) OnRESPMOD(request []byte) (response []byte, err error) {
	return request, nil
}

func Test_net_provider(t *testing.T) {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly})
	w := &Wrapper{
		Provider: server.ProviderNet,
		Handler:  &DummyICAPHandler{},
		Log:      logger,
	}

	var port = 8080
	var url = fmt.Sprintf("http://localhost:%d", port)

	go w.Start(port, context.Background())

	t.Run("it should return 200 when REQMOD is sent", func(t *testing.T) {
		req, err := http.NewRequest("REQMOD", url, nil)
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("it should return 200 when RESPMOD is sent", func(t *testing.T) {
		req, err := http.NewRequest("RESPMOD", url, nil)
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		assert.Equal(t, 200, resp.StatusCode)
	})
}

func Test_gnet_provider(t *testing.T) {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly})
	w := &Wrapper{
		Provider: server.ProviderGNet,
		Handler:  &DummyICAPHandler{},
		Log:      logger,
	}

	var port = 8081
	var url = fmt.Sprintf("http://localhost:%d", port)

	go w.Start(port, context.Background())

	t.Run("it should return 200 when REQMOD is sent", func(t *testing.T) {
		req, err := http.NewRequest("REQMOD", url, nil)
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("it should return 200 when RESPMOD is sent", func(t *testing.T) {
		req, err := http.NewRequest("RESPMOD", url, nil)
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		assert.Equal(t, 200, resp.StatusCode)
	})
}
