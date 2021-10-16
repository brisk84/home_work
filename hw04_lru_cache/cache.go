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
	var res bool
	if l.items[key] == nil {
		item := new(cacheItem)
		item.key = key
		item.value = value

		listItem := l.queue.PushFront(item)
		l.items[key] = listItem
		res = false
	} else {
		l.items[key].Value.(*cacheItem).value = value
		l.queue.MoveToFront(l.items[key])
		res = true
	}

	if l.queue.Len() > l.capacity {
		listItem := l.queue.Back()
		item := listItem.Value.(*cacheItem)
		l.queue.Remove(listItem)
		delete(l.items, item.key)
		res = false
	}

	return res
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
	l.capacity = 0
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
