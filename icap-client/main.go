package main

import (
	"bytes"
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	var host string
	flag.StringVar(&host, "host", "http://127.0.0.1:8080", "ICAP endpoint")
	flag.Parse()

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly})

	latencies := make([]time.Duration, 0, 1000)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			latency := measureLatency(host)
			mu.Lock()
			latencies = append(latencies, latency)
			mu.Unlock()
		}()
	}

	wg.Wait()

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	percentile95 := latencies[int(float64(len(latencies))*0.95)]
	percentile99 := latencies[int(float64(len(latencies))*0.99)]

	logger.Info().Msgf("95th percentile latency: %s", percentile95.String())
	logger.Info().Msgf("99th percentile latency: %s", percentile99.String())
}

func measureLatency(url string) time.Duration {
	body := bytes.NewBuffer([]byte(`{
		"title": "Post title",
		"body": "Post description",
		"userId": 1
	}`))

	start := time.Now()
	req, _ := http.NewRequest("REQMOD", url, body)
	http.DefaultClient.Do(req)
	return time.Since(start)
}
