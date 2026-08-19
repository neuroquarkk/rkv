package dispatcher

import (
	"context"
	pb "rkv/gen"
	"rkv/internal/constants"
)

func (c *Client) Put(
	ctx context.Context,
	key string,
	data []byte,
) error {
	client, err := c.getClient(key)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, constants.PutTimeout)
	defer cancel()

	_, err = client.Put(ctx, &pb.PutRequest{
		Key:  key,
		Data: data,
	})
	return err
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	client, err := c.getClient(key)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, constants.GetTimeout)
	defer cancel()

	resp, err := client.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		return nil, err
	}

	data := resp.Data
	return data, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	client, err := c.getClient(key)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DeleteTimeout)
	defer cancel()

	_, err = client.Delete(ctx, &pb.DeleteRequest{Key: key})
	return err
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	client, err := c.getClient(key)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, constants.ExistsTimeout)
	defer cancel()

	resp, err := client.Exists(ctx, &pb.ExistsRequest{Key: key})
	if err != nil {
		return false, err
	}
	return resp.Exists, nil
}

func (c *Client) Info(ctx context.Context) ([]string, error) {
	addr := make([]string, 0, len(c.clients))

	ctx, cancel := context.WithTimeout(ctx, constants.InfoTimeout)
	defer cancel()

	c.mu.RLock()
	clients := make([]pb.ShardServiceClient, 0, len(c.clients))
	for _, sc := range c.clients {
		clients = append(clients, sc)
	}
	c.mu.RUnlock()

	for _, sc := range clients {
		resp, err := sc.Info(ctx, &pb.InfoRequest{})
		if err != nil {
			return nil, err
		}
		addr = append(addr, resp.Name)
	}
	return addr, nil
}
