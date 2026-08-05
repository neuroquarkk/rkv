package client

import (
	"context"
	"fmt"
	pb "rkv/gen"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Client struct {
	conn   *grpc.ClientConn
	Client pb.ShardServiceClient
}

func New(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	healthClient := healthpb.NewHealthClient(conn)
	resp, err := healthClient.Check(timeoutCtx, &healthpb.HealthCheckRequest{})
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping shard server: %w", err)
	}

	if resp.Status != healthpb.HealthCheckResponse_SERVING {
		conn.Close()
		return nil, fmt.Errorf("shard server not ready... status: %s", resp.Status)
	}

	client := pb.NewShardServiceClient(conn)
	return &Client{conn, client}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Put(ctx context.Context, key string, data []byte) error {
	_, err := c.Client.Put(ctx, &pb.PutRequest{Key: key, Data: data})
	if err != nil {
		return err
	}
	return nil
}
