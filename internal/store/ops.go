package store

import (
	"bytes"
	"time"
)

func (s *Store) Put(key string, value []byte) {
	b := s.getBucket(key)
	b.mu.Lock()
	defer b.mu.Unlock()

	// deep copy incoming byte slice to prevent caller mutations
	newVal := bytes.Clone(value)
	e := &entry{value: newVal}
	e.accessedAt.Store(time.Now().UnixNano())

	if old, ok := b.data[key]; ok {
		b.size -= uint64(len(old.value))
	}

	b.size += uint64(len(newVal))
	b.data[key] = e

	b.evictLocked()
}

func (s *Store) Get(key string) ([]byte, error) {
	b := s.getBucket(key)
	b.mu.RLock()
	defer b.mu.RUnlock()

	if entry, ok := b.data[key]; ok {
		// deep copy outgoing data to prevent accidental caller mutations
		entry.accessedAt.Store(time.Now().UnixNano())
		return bytes.Clone(entry.value), nil
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

	e, ok := b.data[key]
	if !ok {
		return ErrKeyNotFound
	}

	b.size -= uint64(len(e.value))
	delete(b.data, key)

	return nil
}
