package net

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"igops.me/icap/server"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Log        zerolog.Logger
	Handler    server.Handler
	httpServer *http.Server
}

func (s Server) ListenAndHandle(port int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.Handle("/", &HTTPHandler{ICAPHandler: s.Handler})

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	go func() {
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.Log.Error().Err(err).Msg("Server startup error")
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := s.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			s.Log.Error().Err(err).Msg("HTTP server shutdown error")
		}
	}
}
