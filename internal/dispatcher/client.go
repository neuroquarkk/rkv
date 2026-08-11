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
	idx := c.pickNode()

	_, err := c.cluster.Clients[idx].Put(ctx, &pb.PutRequest{
		Key:  key,
		Data: data,
	})
	return err
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	idx := c.pickNode()

	resp, err := c.cluster.Clients[idx].Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		return nil, err
	}

	data := resp.Data
	return data, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	idx := c.pickNode()

	_, err := c.cluster.Clients[idx].Delete(ctx, &pb.DeleteRequest{Key: key})
	return err
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	idx := c.pickNode()
	resp, err := c.cluster.Clients[idx].Exists(ctx, &pb.ExistsRequest{Key: key})
	if err != nil {
		return false, err
	}
	return resp.Exists, nil
}

func (c *Client) Info(ctx context.Context) (string, error) {
	idx := c.pickNode()

	resp, err := c.cluster.Clients[idx].Info(ctx, &pb.InfoRequest{})
	if err != nil {
		return "", err
	}
	return resp.Name, nil
}
