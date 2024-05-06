package main

import (
	"context"
	"flag"
	"github.com/rs/zerolog"
	"igops.me/icap/server"
	"igops.me/icap/wrapper"
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
}

func main() {
	var port int
	var providerStr string

	flag.StringVar(&providerStr, "provider", "net", "server provider")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	var provider = parseProvider(providerStr)
	if provider == server.ProviderUnknown {
		log.Fatal().Msgf("Unknown provider: %s", providerStr)
	}

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly})
	w := &wrapper.Wrapper{
		Provider: provider,
		Handler:  &CustomICAPHandler{},
		Log:      logger,
	}

	w.Start(port, context.Background())
}

func parseProvider(provider string) server.Provider {
	switch provider {
	case "gnet":
		return server.ProviderGNet
	case "net":
		return server.ProviderNet
	default:
		return server.ProviderUnknown
	}
}

type CustomICAPHandler struct {
}

func (h *CustomICAPHandler) OnREQMOD(request []byte) (response []byte, err error) {
	return request, nil
}

func (h *CustomICAPHandler) OnRESPMOD(request []byte) (response []byte, err error) {
	return request, nil
}
