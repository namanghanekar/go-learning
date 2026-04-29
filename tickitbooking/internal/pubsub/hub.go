package pubsub

import "sync"

type Hub[T any] struct {
	mu          sync.RWMutex
	subscribers map[chan T]struct{}
}

func NewHub[T any]() *Hub[T] {
	return &Hub[T]{subscribers: map[chan T]struct{}{}}
}

func (h *Hub[T]) Subscribe(buffer int) chan T {
	ch := make(chan T, buffer)
	h.mu.Lock()
	h.subscribers[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub[T]) Unsubscribe(ch chan T) {
	h.mu.Lock()
	if _, ok := h.subscribers[ch]; ok {
		delete(h.subscribers, ch)
		close(ch)
	}
	h.mu.Unlock()
}

func (h *Hub[T]) Publish(msg T) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subscribers {
		select {
		case ch <- msg:
		default:
		}
	}
}
