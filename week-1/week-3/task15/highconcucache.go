package main

import (
	"fmt"
	"sync"
)

type MapCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMapCache() *MapCache {
	return &MapCache{
		data: make(map[string]string),
	}
}
func (m *MapCache) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}
func (m *MapCache) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}
func main() {
	cache := NewMapCache()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			cache.Set(key, "value")
			fmt.Println("SET:", key)
		}(i)
	}
	wg.Wait()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			val, ok := cache.Get(key)
			fmt.Println("GET:", key, val, ok)
		}(i)
	}
	wg.Wait()
}
