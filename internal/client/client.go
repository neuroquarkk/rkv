package client

import (
	pb "rkv/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	Client pb.ShardServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewShardServiceClient(conn)
	return &Client{conn, client}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
