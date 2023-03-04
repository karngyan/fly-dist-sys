package main

import "sync"

type store struct {
	mu sync.RWMutex
	m  map[int64]bool

	neighbors []string
}

func newStore() *store {
	return &store{
		m:         map[int64]bool{},
		neighbors: []string{},
	}
}

func (s *store) addM(v int64) {
	s.mu.Lock()
	s.m[v] = true
	s.mu.Unlock()
}

func (s *store) getMByKey(k int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[k]
}

func (s *store) getM() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var r []int64
	for k := range s.m {
		r = append(r, k)
	}
	return r
}

func (s *store) setNeighbors(v []string) {
	s.neighbors = v
}

func (s *store) getNeighbors() []string {
	return s.neighbors
}
