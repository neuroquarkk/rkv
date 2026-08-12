package poller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rkv/internal/registry"
	"time"
)

var (
	client    *http.Client
	targetUrl string
)

func Start(ctx context.Context, url string, ch chan []string) {
	targetUrl = "http://" + url + "/members"

	client = &http.Client{
		Timeout: 2 * time.Second,
	}

	var initialAddrs []string

	// Synchronous boot barrier
	// block router from starting the server
	// untill registry provides initial cluster
	for {
		addrs, err := fetchMembers()
		if err == nil {
			initialAddrs = addrs
			break
		}
		log.Println(err)

		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
		}
	}

	ch <- initialAddrs

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				addrs, err := fetchMembers()
				if err != nil {
					log.Println(err)
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

func fetchMembers() ([]string, error) {
	resp, err := client.Get(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get members list: %w\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var respData registry.MembersResp
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w\n", err)
	}

	return respData.Members, nil
}
