package reporter

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rkv/internal/registry"
	"time"
)

var (
	client    *http.Client
	targetUrl string
)

func Start(ctx context.Context, url string, addr string, d time.Duration) {
	client = &http.Client{
		Timeout: 2 * time.Second,
	}
	targetUrl = "http://" + url + "/heartbeat"

	// discover immediately with a soft dependency
	// periodic retries handle failures in the background
	// allowing the member to start without waiting for discovery to succeed
	do(addr)

	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				do(addr)
			}
		}
	}()
}

func do(addr string) {
	req := &registry.HeartbeatReq{
		Address: addr,
	}

	data, err := json.Marshal(&req)
	if err != nil {
		log.Printf("failed to marshal: %v\n", err)
		return
	}

	resp, err := client.Post(
		targetUrl,
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		log.Printf("failed to send heartbeat: %v\n", err)
		return
	}
	if resp.StatusCode != http.StatusAccepted {
		log.Printf("registry rejected heartbeat %d\n", resp.StatusCode)
	}
	resp.Body.Close()
}
