package registry

import (
	"context"
	"errors"
	"fmt"
	pb "rkv/gen"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Registry struct {
	conns   []*grpc.ClientConn
	Clients map[string]pb.ShardServiceClient
}

type connResult struct {
	addr   string
	conn   *grpc.ClientConn
	client pb.ShardServiceClient
	err    error
}

func New(
	ctx context.Context,
	addrs string,
) (*Registry, error) {
	members := strings.Split(addrs, ",")
	if len(members) == 0 || addrs == "" {
		return nil, errors.New("no member addresses provided")
	}

	conns := make([]*grpc.ClientConn, 0, len(members))
	clients := make(map[string]pb.ShardServiceClient, len(members))

	var wg sync.WaitGroup
	ch := make(chan connResult, len(members))

	for _, m := range members {
		wg.Go(func() {
			conn(ctx, m, ch)
		})
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var errs []error
	for res := range ch {
		if res.err != nil {
			errs = append(errs, res.err)
			continue
		}
		conns = append(conns, res.conn)
		clients[res.addr] = res.client
	}

	if len(errs) > 0 {
		for _, c := range conns {
			c.Close()
		}
		joinedErr := errors.Join(errs...)

		return nil, fmt.Errorf(
			"registry initalization failed (%d/%d members unreachable):\n%w",
			len(errs),
			len(members),
			joinedErr,
		)
	}

	return &Registry{
		conns:   conns,
		Clients: clients,
	}, nil
}

func (r *Registry) Close() {
	for _, c := range r.conns {
		c.Close()
	}
}

func conn(ctx context.Context, addr string, ch chan<- connResult) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		err := fmt.Errorf("failed to create client for %s: %w", addr, err)
		ch <- connResult{err: err}
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	healthClient := healthpb.NewHealthClient(conn)
	resp, err := healthClient.Check(timeoutCtx, &healthpb.HealthCheckRequest{})

	if err != nil || resp.Status != healthpb.HealthCheckResponse_SERVING {
		conn.Close()
		err := fmt.Errorf("failed to ping shard %s: %w", addr, err)
		ch <- connResult{err: err}
		return
	}

	client := pb.NewShardServiceClient(conn)
	ch <- connResult{addr: addr, conn: conn, client: client}
}
