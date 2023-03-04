package main

import "sync"

type store struct {
	mu sync.RWMutex
	m  []int64

	tu sync.RWMutex
	t  map[string][]string
}

func newStore() *store {
	return &store{
		m: []int64{},
		t: map[string][]string{},
	}
}

func (s *store) addM(v int64) {
	s.mu.Lock()
	s.m = append(s.m, v)
	s.mu.Unlock()
}

func (s *store) getM() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m
}

func (s *store) addT(k string, v []string) {
	s.tu.Lock()
	s.t[k] = v
	s.tu.Unlock()
}

func (s *store) getT() map[string][]string {
	s.tu.RLock()
	defer s.tu.RUnlock()
	return s.t
}
