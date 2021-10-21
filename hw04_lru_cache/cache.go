package hw04lrucache

import (
    "fmt"
    "sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

    mu sync.Mutex
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
        item.key = string(key)      // ???
        item.value = value

        listItem := l.queue.PushFront(item)
        l.items[key] = listItem
        res = false
    } else {
        //fmt.Println(l.queue.Back().Value)
        //listItem := l.items[key]
        l.items[key].Value.(*cacheItem).value = value
        l.queue.MoveToFront(l.items[key])
        //fmt.Println(l.queue.Back().Value)
        //l.queue.MoveToFront(listItem)
        res = true
    }

    if l.queue.Len() > l.capacity {
        listItem := l.queue.Back()
        if listItem == nil {
            fmt.Printf("Len: %d, cap: %d\n", l.queue.Len(), l.capacity)
            fmt.Println(l.queue.Back())
            fmt.Println(l.queue.Front())
            fmt.Println(l.queue.Front())
            fmt.Println(l.queue.Front())
            fmt.Println(l.queue.Front())
            //return false
        }
        item := listItem.Value.(*cacheItem)
        //fmt.Printf("Removing %s\n", item.key)
        //l.queue.Remove(listItem)
        l.queue.Remove(l.queue.Back())
        delete(l.items, Key(item.key))

        //l.queue.Remove(l.queue.Back())

        res = false
    }

    //fmt.Printf("Len: %d, Cap: %d\n", l.queue.Len(), l.capacity)
    //fmt.Println()
    return res
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
    l.mu.Lock()
    defer l.mu.Unlock()
    if l.items[key] != nil {
        val := l.items[key].Value
        //fmt.Printf("Key: %s, val: %d\n", key, val)
        //listItem := l.items[key]
        //l.queue.MoveToFront(listItem)
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
    l.items = make(map[Key]*ListItem, 0)
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
