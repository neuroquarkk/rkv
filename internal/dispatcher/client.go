package dispatcher

import (
	"context"
	pb "rkv/gen"
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

	_, err = client.Delete(ctx, &pb.DeleteRequest{Key: key})
	return err
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	client, err := c.getClient(key)
	if err != nil {
		return false, err
	}

	resp, err := client.Exists(ctx, &pb.ExistsRequest{Key: key})
	if err != nil {
		return false, err
	}
	return resp.Exists, nil
}

func (c *Client) Info(ctx context.Context) ([]string, error) {
	addr := make([]string, 0, len(c.clients))

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, sc := range c.clients {
		resp, err := sc.Info(ctx, &pb.InfoRequest{})
		if err != nil {
			return nil, err
		}
		addr = append(addr, resp.Name)
	}
	return addr, nil
}
