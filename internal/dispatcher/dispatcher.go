package dispatcher

import (
	"math/rand/v2"
	"rkv/internal/registry"
)

type Client struct {
	cluster *registry.Registry
}

func New(r *registry.Registry) *Client {
	return &Client{cluster: r}
}

// TODO: replace random selection with deterministic routing
func (c *Client) pickNode() int {
	return rand.IntN(len(c.cluster.Clients))
}
