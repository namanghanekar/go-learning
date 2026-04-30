package graph

import (
	"graphql-service/client"
	"graphql-service/graph/model"
	"sync"
)

type Resolver struct {
	Client *client.Client

	mu          sync.RWMutex
	subscribers map[int]chan *model.Seat
	nextID      int
}

func (r *Resolver) AddSubscriber() (int, chan *model.Seat) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.subscribers == nil {
		r.subscribers = make(map[int]chan *model.Seat)
	}

	r.nextID++
	id := r.nextID
	ch := make(chan *model.Seat, 20)
	r.subscribers[id] = ch
	return id, ch
}

func (r *Resolver) RemoveSubscriber(id int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.subscribers[id]; ok {
		delete(r.subscribers, id)
		close(ch)
	}
}
func (r *Resolver) PublishSeatUpdate(seat *model.Seat) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, ch := range r.subscribers {
		select {
		case ch <- seat:
		default:
		}
	}
}
