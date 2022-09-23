package sync

import (
	"sync"
)

type WaitMapObject struct {
	wg   map[string]bool
	mu   sync.Mutex
	cond sync.Cond
}

func WaitMap() *WaitMapObject {
	m := &WaitMapObject{}
	m.wg = make(map[string]bool)
	m.cond.L = &m.mu
	return m
}

func (m *WaitMapObject) Wait(name string) {
	m.mu.Lock()
	for m.wg[name] {
		m.cond.Wait()
	}
	m.mu.Unlock()
}

func (m *WaitMapObject) UnLock(name string) {
	m.mu.Lock()
	m.wg[name] = false
	m.mu.Unlock()
	m.cond.Broadcast()
}

func (m *WaitMapObject) Lock(name string) {
	m.Wait(name)
	m.mu.Lock()
	m.wg[name] = true
	m.mu.Unlock()
}
