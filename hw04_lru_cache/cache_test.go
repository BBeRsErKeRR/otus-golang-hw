package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		lru := NewCache(2)
		lru.Set("a", 1)
		lru.Set("b", 2)
		lru.Set("c", 3)

		val, ok := lru.Get("a")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge old by set", func(t *testing.T) {
		lru := NewCache(3)
		lru.Set("a", 1)  // [a:1]
		lru.Set("b", 2)  // [a:1, b:2]
		lru.Set("a", 3)  // [b:3, a:1]
		lru.Set("c", 4)  // [c:4, b:3, a:1]
		lru.Set("a", 5)  // [a:5, c:4, b:3]
		lru.Set("d", 10) // [d:10, a:5, c:4]

		val, ok := lru.Get("b")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge old by get", func(t *testing.T) {
		lru := NewCache(3)
		lru.Set("a", 1) // [a:1]
		lru.Set("b", 2) // [b:2, a:1]
		lru.Set("c", 4) // [c:4, b:3, a:1]
		lru.Get("a")    // [a:1, c:4, b:3]
		lru.Get("c")    // [c:4, a:1, b:3]
		lru.Get("b")    // [b:3, c:4, a:1]
		lru.Get("a")    // [a:1, b:3, c:4]
		lru.Set("d", 0) // [d:0, a:1, b:3]

		val, ok := lru.Get("c")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear cache", func(t *testing.T) {
		lru := NewCache(3)
		lru.Set("a", 1) // [a:1]
		lru.Set("b", 2) // [b:2, a:1]
		lru.Set("c", 4) // [c:4, b:3, a:1]
		lru.Clear()
		for _, v := range [...]Key{"a", "b", "c"} {
			val, ok := lru.Get(v)
			require.False(t, ok)
			require.Nil(t, val)
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
	c.Clear()
}
