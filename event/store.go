package event

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const Size = 1024

type MemoryStore struct {
	Epoch    uint64
	Counter  atomic.Uint64
	Queue    chan Event
	Emitter  Emitter
	Ingester Ingester
	SSEEmit  bool
	sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	m := &MemoryStore{
		Epoch: uint64(time.Now().UnixNano()),
		Queue: make(chan Event, Size),
	}
	m.Emitter = Emitter{m}
	m.Ingester = Ingester{m}
	for i := 0; i < 4; i++ {
		msg := Payload{
			"msg": fmt.Sprintf("hello %d", i),
		}
		m.DispatchEvent(msg)
	}
	return m
}

func (m *MemoryStore) NewEvent(data Payload) Event {
	e := Event{
		Data: data,
		Is:   true,
		Id:   m.Counter.Add(1),
	}
	return e
}

func (m *MemoryStore) DispatchEvent(data Payload) {
	e := m.NewEvent(data)
	m.Enqueue(e)
}

func (m *MemoryStore) Enqueue(e Event) uint64 {
	m.Queue <- e
	return e.Id
}
