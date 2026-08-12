package dispatcher

import (
	"errors"
	"hash/fnv"
	pb "rkv/gen"
)

var ErrNoMembers = errors.New("no members available")

func (c *Client) getClient(key string) (pb.ShardServiceClient, error) {
	var hashVal uint64
	var target string

	f := fnv.New64a()

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.clients) == 0 {
		return nil, ErrNoMembers
	}

	// Rendezvous Hashing
	// hash the key along with the active member's address
	// member which produce highest hash wins
	for addr := range c.clients {
		f.Write([]byte(key))
		f.Write([]byte(":"))
		f.Write([]byte(addr))
		s := f.Sum64()

		if s > hashVal {
			hashVal = s
			target = addr
		}

		f.Reset()
	}

	return c.clients[target], nil
}
