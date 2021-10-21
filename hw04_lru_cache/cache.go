package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.capacity == 0 {
		return false
	}

	if l.queue.Len() >= l.capacity {
		listItem := l.queue.Back()
		item := listItem.Value.(*cacheItem)
		l.queue.Remove(listItem)
		delete(l.items, item.key)
	}

	item, ok := l.items[key]
	if ok {
		item.Value.(*cacheItem).value = value
		l.queue.MoveToFront(item)
		return true
	}

	cItem := &cacheItem{key, value}
	listItem := l.queue.PushFront(cItem)
	l.items[key] = listItem
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.items[key] != nil {
		val := l.items[key].Value
		l.queue.MoveToFront(l.items[key])
		return val.(*cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
