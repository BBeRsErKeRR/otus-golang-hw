package hw04lrucache

import "sync"

var syncMutex = sync.Mutex{}

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
}

type ListValue struct {
	key   Key
	value interface{}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	isExist := false
	var listValue ListValue
	listValue.key = key
	listValue.value = value

	syncMutex.Lock()
	item, ok := lru.items[key]
	syncMutex.Unlock()

	if ok {
		item.Value = listValue
		lru.queue.MoveToFront(item)
		isExist = true
	} else {
		syncMutex.Lock()
		if lru.capacity == lru.queue.Len() {
			listitem := lru.queue.Back()
			lru.queue.Remove(listitem)
			delete(lru.items, listitem.Value.(ListValue).key)
		}
		lru.items[key] = lru.queue.PushFront(listValue)
		syncMutex.Unlock()
	}

	return isExist
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	syncMutex.Lock()
	item, ok := lru.items[key]
	syncMutex.Unlock()
	if ok {
		lru.queue.MoveToFront(item)
		return item.Value.(ListValue).value, true
	}
	return nil, false
}

func (lru *lruCache) Clear() {
	syncMutex.Lock()
	lru.items = make(map[Key]*ListItem)
	lru.queue = NewList()
	syncMutex.Unlock()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
