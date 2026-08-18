package store

import (
	"bytes"
	"errors"
	"hash/maphash"
	"sync"

	"rkv/internal/constants"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

const numBuckets = uint64(constants.NumBuckets)

type bucket struct {
	mu   sync.RWMutex
	data map[string][]byte
}

type Store struct {
	buckets [numBuckets]*bucket
	seed    maphash.Seed
}

func New() *Store {
	var buckets [numBuckets]*bucket

	for i := range numBuckets {
		buckets[i] = &bucket{
			data: make(map[string][]byte),
		}
	}

	return &Store{
		buckets: buckets,
		seed:    maphash.MakeSeed(),
	}
}

func (s *Store) Put(key string, value []byte) {
	b := s.getBucket(key)
	b.mu.Lock()
	defer b.mu.Unlock()

	// deep copy incoming byte slice to prevent caller mutations
	b.data[key] = bytes.Clone(value)
}

func (s *Store) Get(key string) ([]byte, error) {
	b := s.getBucket(key)
	b.mu.RLock()
	defer b.mu.RUnlock()

	if data, ok := b.data[key]; ok {
		// deep copy outgoing data to prevent accidental caller mutations
		return bytes.Clone(data), nil
	}

	return nil, ErrKeyNotFound
}

func (s *Store) Exists(key string) bool {
	b := s.getBucket(key)
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, ok := b.data[key]

	return ok
}

func (s *Store) Delete(key string) error {
	b := s.getBucket(key)
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.data[key]; !ok {
		return ErrKeyNotFound
	}

	delete(b.data, key)
	return nil
}

func (s *Store) getBucket(key string) *bucket {
	h := maphash.String(s.seed, key)
	i := h & (numBuckets - 1)
	return s.buckets[i]
}
