package lru

import (
	"container/list"
	"sync"
)

type Item struct {
	Key   string
	Value string
}

type LRUCache struct {
	sync.RWMutex
	capacity int
	items    map[string]*list.Element
	queue    *list.List
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (cache *LRUCache) Set(key, value string) bool {
	cache.Lock()
	defer cache.Unlock()

	if element, exists := cache.items[key]; exists {
		cache.queue.MoveToFront(element)

		element.Value.(*Item).Value = value

		return false
	}

	if cache.queue.Len() == cache.capacity {
		cache.purge()
	}

	item := &Item{
		Key:   key,
		Value: value,
	}

	element := cache.queue.PushFront(item)
	cache.items[key] = element

	return true
}

func (cache *LRUCache) Get(key string) (string, bool) {
	cache.RLock()
	defer cache.RUnlock()

	element, exists := cache.items[key]
	if !exists {
		return "", false
	}

	cache.queue.MoveToFront(element)

	return element.Value.(*Item).Value, true
}

func (cache *LRUCache) Delete(key string) bool {
	cache.Lock()
	defer cache.Unlock()

	element, exists := cache.items[key]
	if !exists {
		return false
	}

	cache.queue.MoveToBack(element)
	cache.purge()

	return true
}

func (cache *LRUCache) purge() {
	if element := cache.queue.Back(); element != nil {
		item := cache.queue.Remove(element).(*Item)
		delete(cache.items, item.Key)
	}
}
