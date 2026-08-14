package reporter

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rkv/internal/constants"
	"rkv/internal/registry"
	"time"
)

type Reporter struct {
	client      *http.Client
	registryUrl string
	addr        string
}

func New(url string, addr string) *Reporter {
	client := &http.Client{
		Timeout: constants.ClientTimeout,
	}

	return &Reporter{
		client:      client,
		registryUrl: url,
		addr:        addr,
	}
}

func (r *Reporter) Start(ctx context.Context, d time.Duration) {
	targetUrl := "http://" + r.registryUrl + "/heartbeat"

	// discover immediately with a soft dependency
	// periodic retries handle failures in the background
	// allowing the member to start without waiting for discovery to succeed
	r.do(targetUrl, http.StatusAccepted)

	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				r.do(targetUrl, http.StatusAccepted)
			}
		}
	}()
}

func (r *Reporter) Leave() {
	targetUrl := "http://" + r.registryUrl + "/leave"
	r.do(targetUrl, http.StatusNoContent)
}

func (r *Reporter) do(targetUrl string, expectedCode int) {
	req := &registry.MemberReq{
		Address: r.addr,
	}

	data, err := json.Marshal(&req)
	if err != nil {
		log.Printf("[SHARD] failed to marshal: %v\n", err)
		return
	}

	resp, err := r.client.Post(
		targetUrl,
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		log.Printf("[SHARD] failed to send: %v\n", err)
		return
	}
	if resp.StatusCode != expectedCode {
		log.Printf("[SHARD] registry rejected %d\n", resp.StatusCode)
	}
	resp.Body.Close()
}
