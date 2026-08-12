package dispatcher

import (
	"rkv/internal/registry"
)

type Client struct {
	cluster *registry.Registry
}

func New(r *registry.Registry) *Client {
	return &Client{cluster: r}
}
