package server

import (
	"context"
	"sync"
)

type Server interface {
	ListenAndHandle(port int, ctx context.Context, wg *sync.WaitGroup)
}
