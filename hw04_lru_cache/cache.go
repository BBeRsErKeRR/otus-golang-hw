package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mx       sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	isExist := false
	var listValue cacheItem
	listValue.key = key
	listValue.value = value

	lru.mx.Lock()
	defer lru.mx.Unlock()
	item, ok := lru.items[key]

	if ok {
		item.Value = listValue
		lru.queue.MoveToFront(item)
		isExist = true
	} else {
		if lru.capacity == lru.queue.Len() {
			listitem := lru.queue.Back()
			lru.queue.Remove(listitem)
			delete(lru.items, listitem.Value.(cacheItem).key)
		}
		lru.items[key] = lru.queue.PushFront(listValue)
	}
	return isExist
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mx.Lock()
	defer lru.mx.Unlock()
	item, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (lru *lruCache) Clear() {
	lru.mx.Lock()
	defer lru.mx.Unlock()
	lru.items = make(map[Key]*ListItem)
	lru.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
