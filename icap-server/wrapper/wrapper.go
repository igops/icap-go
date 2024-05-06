package wrapper

import (
	"context"
	"github.com/rs/zerolog"
	"igops.me/icap/server"
	"igops.me/icap/server/providers/gnet"
	"igops.me/icap/server/providers/net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Wrapper struct {
	Provider server.Provider
	Handler  server.Handler
	Log      zerolog.Logger
	server   server.Server
}

func (w *Wrapper) Start(port int, ctx context.Context) {
	if w.server != nil {
		w.Log.Warn().Msg("Server already started")
		return
	}

	var s server.Server
	switch w.Provider {
	case server.ProviderNet:
		s = &net.Server{
			Handler: w.Handler,
			Log:     w.Log,
		}
	case server.ProviderGNet:
		s = &gnet.Server{
			Handler:   w.Handler,
			Log:       w.Log,
			Multicore: true,
		}
	default:
		w.Log.Fatal().Msg("Should never happen")
	}

	cancelCtx, shutdown := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go s.ListenAndHandle(port, cancelCtx, &wg)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	w.Log.Info().Msg("Gracefully shutting down...")
	shutdown()

	wg.Wait()
	w.Log.Info().Msg("Shutdown complete")
}
