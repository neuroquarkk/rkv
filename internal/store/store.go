package store

import (
	"errors"
	"hash/maphash"
	"sync"
	"sync/atomic"

	"rkv/internal/constants"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

const numBuckets = uint64(constants.NumBuckets)

type entry struct {
	value      []byte
	accessedAt atomic.Int64 // atomic because Get() mutates this under RLock
}

type bucket struct {
	mu     sync.RWMutex
	data   map[string]*entry
	size   uint64
	budget uint64
}

type Store struct {
	buckets [numBuckets]*bucket
	seed    maphash.Seed
}

func New() *Store {
	var buckets [numBuckets]*bucket
	budget := constants.MaxStorageSize / numBuckets

	for i := range numBuckets {
		buckets[i] = &bucket{
			data:   make(map[string]*entry),
			budget: budget,
		}
	}

	return &Store{
		buckets: buckets,
		seed:    maphash.MakeSeed(),
	}
}

func (s *Store) getBucket(key string) *bucket {
	h := maphash.String(s.seed, key)
	i := h & (numBuckets - 1)
	return s.buckets[i]
}
