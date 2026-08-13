package server

import (
	"context"
	"log"
	"sync"
	"time"
)

type Server struct {
	mu            sync.RWMutex
	state         map[string]time.Time
	staleInterval time.Duration
	tag           uint64
}

func New(staleInterval time.Duration) *Server {
	return &Server{
		state:         make(map[string]time.Time),
		staleInterval: staleInterval,
		tag:           0,
	}
}

func (s *Server) Sweeper(ctx context.Context, d time.Duration) {
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.mu.Lock()
				// Background garbage collection
				// evict members that have missed their heartbeat window
				for addr, t := range s.state {
					if time.Since(t) > s.staleInterval {
						log.Printf("[REGISTRY] removing stale member: %s\n", addr)
						delete(s.state, addr)
						s.tag++
					}
				}
				s.mu.Unlock()
			}
		}
	}()
}
