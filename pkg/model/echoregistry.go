package model

import "sync"

type EchoRegistry struct {
	sessions sync.Map
}

func NewEchoRegistry() *EchoRegistry {
	return &EchoRegistry{
		sessions: sync.Map{},
	}
}

func (e *EchoRegistry) Register(sessionID string, ch chan string) {
	e.sessions.Store(sessionID, ch)
}

func (e *EchoRegistry) Unregister(sessionID string) {
	e.sessions.Delete(sessionID)
}

func (e *EchoRegistry) Send(sessionID string, msg string) {
	if v, ok := e.sessions.Load(sessionID); ok {
		ch := v.(chan string)
		ch <- msg
	}
}
