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

	if item, ok := l.items[key]; ok {
		l.items[key].Value = value
		l.queue.MoveToFront(item)

		return true
	}

	newItem := l.queue.PushFront(value)
	l.items[key] = newItem

	if l.queue.Len() > l.capacity {
		lastItem := l.queue.Back()
		l.queue.Remove(lastItem)

		for key, item := range l.items {
			if item == lastItem {
				delete(l.items, key)
				break
			}
		}
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	item, ok := l.items[key]

	if !ok {
		return nil, false
	}

	l.queue.MoveToFront(item)

	return item.Value, true
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
