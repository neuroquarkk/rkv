package store

import (
	"bytes"
	"errors"
	"sync"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Store struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func New() *Store {
	return &Store{
		data: make(map[string][]byte),
	}
}

func (s *Store) Put(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// deep copy incoming byte slice to prevent caller mutations
	s.data[key] = bytes.Clone(value)
}

func (s *Store) Get(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if data, ok := s.data[key]; ok {
		// deep copy outgoing data to prevent accidental caller mutations
		return bytes.Clone(data), nil
	}

	return nil, ErrKeyNotFound
}

func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[key]

	return ok
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; !ok {
		return ErrKeyNotFound
	}

	delete(s.data, key)
	return nil
}
