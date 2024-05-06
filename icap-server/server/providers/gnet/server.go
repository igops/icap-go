package gnet

import (
	"context"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/rs/zerolog"
	"igops.me/icap/server"
	"sync"
	"time"
)

type Server struct {
	Log       zerolog.Logger
	Handler   server.Handler
	Multicore bool
}

func (s Server) ListenAndHandle(port int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	addr := fmt.Sprintf("tcp://127.0.0.1:%d", port)

	go func() {
		listener := &Listener{handler: s.Handler, log: s.Log}
		err := gnet.Serve(listener, addr, gnet.WithMulticore(s.Multicore))
		if err != nil {
			s.Log.Error().Err(err).Msg("Server startup error")
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := gnet.Stop(shutdownCtx, addr)
		if err != nil {
			s.Log.Error().Err(err).Msg("HTTP server shutdown error")
		}
	}
}
