package adapters

import (
	"sync"

	"github.com/cmilhench/kvs/internal/ports"
)

var _ ports.Store = (*MemoryStore)(nil)

type MemoryStore struct {
	mu      sync.Mutex
	storage map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		storage: make(map[string]string),
	}
}

func (s *MemoryStore) Get(key string) (string, error) {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	if value, ok := s.storage[key]; ok {
		return value, nil
	}
	return "", nil
}

func (s *MemoryStore) Put(key string, value string) error {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	s.storage[key] = value
	return nil
}

func (s *MemoryStore) Append(key string, value string) (string, error) {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	if old, ok := s.storage[key]; ok {
		s.storage[key] = old + value
		return old, nil
	}
	s.storage[key] = value
	return "", nil
}
