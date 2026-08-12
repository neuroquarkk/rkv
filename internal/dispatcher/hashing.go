package dispatcher

import (
	"hash/fnv"
	pb "rkv/gen"
)

func (c *Client) getClient(key string) pb.ShardServiceClient {
	var hashVal uint64
	var target string

	f := fnv.New64a()
	for addr := range c.cluster.Clients {
		f.Write([]byte(key))
		f.Write([]byte(addr))
		s := f.Sum64()

		if s > hashVal {
			hashVal = s
			target = addr
		}

		f.Reset()
	}

	return c.cluster.Clients[target]
}
