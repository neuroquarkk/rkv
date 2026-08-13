package poller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rkv/internal/constants"
	"rkv/internal/registry"
	"strconv"
	"time"
)

type Poller struct {
	client    *http.Client
	targetUrl string
	tag       uint64
}

func New(url string) *Poller {
	client := &http.Client{
		Timeout: constants.ClientTimeout,
	}
	targetUrl := "http://" + url + "/members"

	return &Poller{
		client:    client,
		targetUrl: targetUrl,
		tag:       0,
	}
}

func (p *Poller) Start(ctx context.Context, ch chan []string) {
	var initialAddrs []string

	// Synchronous boot barrier
	// block router from starting the server
	// untill registry provides initial cluster
	for {
		addrs, err := p.fetchMembers()
		if err == nil {
			initialAddrs = addrs
			break
		}
		log.Printf("[ROUTER] %v\n", err)

		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
		}
	}

	ch <- initialAddrs

	go func() {
		ticker := time.NewTicker(constants.PollerTick)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				addrs, err := p.fetchMembers()
				if err != nil {
					log.Printf("[ROUTER] %v\n", err)
					continue
				}

				if addrs == nil {
					continue
				}

				// Latest value only
				// if the dispatcher is slow, stale buffered data is discarded
				// so it always processes the newest state
				select {
				case <-ch:
				default:
				}

				ch <- addrs
			}
		}
	}()
}

func (p *Poller) fetchMembers() ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, p.targetUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if p.tag > 0 {
		req.Header.Set("x-tag", strconv.FormatUint(p.tag, 10))
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get members list: %w", err)
	}

	if resp.StatusCode == http.StatusNotModified {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	tag := resp.Header.Get("x-tag")
	if tag != "" {
		parsedTag, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse x-tag: %w", err)
		}

		log.Printf("[ROUTER] tag changed: tag %d -> %d\n", p.tag, parsedTag)
		p.tag = parsedTag
	} else {
		log.Println("[ROUTER] missing x-tag")
		p.tag = 0
	}

	var respData registry.MembersResp
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w\n", err)
	}

	return respData.Members, nil
}
