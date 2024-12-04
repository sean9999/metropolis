package event

import (
	"fmt"
	"sync/atomic"
)

const Size = 1024

type Bus struct {
	Counter atomic.Uint64
	Queue   chan Event
}

func NewBus() *Bus {
	m := &Bus{
		Queue: make(chan Event, Size),
	}
	for i := 0; i < 4; i++ {
		msg := Payload{
			"msg": fmt.Sprintf("hello %d", i),
		}
		m.DispatchEvent(msg)
	}
	return m
}

func (m *Bus) NewEvent(data Payload) Event {
	e := Event{
		Data: data,
		Is:   true,
		Id:   m.Counter.Add(1),
	}
	return e
}

func (m *Bus) DispatchEvent(data Payload) {
	e := m.NewEvent(data)
	m.Enqueue(e)
}

func (m *Bus) Enqueue(e Event) uint64 {
	m.Queue <- e
	return e.Id
}
