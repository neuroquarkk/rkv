package dispatcher

import (
	"context"
	"log"
	pb "rkv/gen"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	mu      sync.RWMutex
	conns   map[string]*grpc.ClientConn
	clients map[string]pb.ShardServiceClient
}

func New() *Client {
	return &Client{
		conns:   make(map[string]*grpc.ClientConn),
		clients: make(map[string]pb.ShardServiceClient),
	}
}

func (c *Client) Start(ctx context.Context, ch <-chan []string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case addrs := <-ch:
				c.process(addrs)
			}
		}
	}()
}

func (c *Client) process(addrs []string) {
	// comparing the absolute state from the registry again our current local state
	newSet := make(map[string]bool, len(addrs))

	c.mu.Lock()
	defer c.mu.Unlock()

	// check for newly registered members and add them to the current state
	for _, addr := range addrs {
		newSet[addr] = true

		if _, exists := c.clients[addr]; !exists {
			if err := c.addMember(addr); err != nil {
				log.Printf("failed to add %s: %v\n", addr, err)
				continue
			}
		}
	}

	// check for dead members and clean them
	for addr := range c.conns {
		if !newSet[addr] {
			c.removeMember(addr)
		}
	}
}

func (c *Client) addMember(addr string) error {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	client := pb.NewShardServiceClient(conn)

	c.conns[addr] = conn
	c.clients[addr] = client
	return nil
}

func (c *Client) removeMember(addr string) {
	c.conns[addr].Close()
	delete(c.conns, addr)
	delete(c.clients, addr)
}
