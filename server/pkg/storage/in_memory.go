package storage

import "grpc-server/pkg/lru"

type KeyValueStorage interface {
	Put(key, value string) bool
	Get(key string) string
	Delete(key string) bool
}

type InMemoryKeyValueStorage struct {
	cache *lru.LRUCache
}

func NewInMemoryKeyValueStorage(capacity int) *InMemoryKeyValueStorage {
	lruCache := lru.NewLRUCache(capacity)

	return &InMemoryKeyValueStorage{
		cache: lruCache,
	}
}

func (s *InMemoryKeyValueStorage) Put(key, value string) bool {
	return s.cache.Set(key, value)
}

func (s *InMemoryKeyValueStorage) Get(key string) (string, bool) {
	return s.cache.Get(key)
}

func (s *InMemoryKeyValueStorage) Delete(key string) bool {
	return s.cache.Delete(key)
}
