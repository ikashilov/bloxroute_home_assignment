package storage

import (
	"assignmentapp/internal/api"
	"errors"
	"sync"
)

var (
	ErrItemAlreadyExists = errors.New("item already exists")
	ErrNoSuchItem        = errors.New("no such item")
)

type Storage struct {
	mu    sync.RWMutex
	cache map[api.ItemKey]api.ItemValue
}

func New(cap int) *Storage {
	var s Storage
	s.cache = make(map[api.ItemKey]api.ItemValue, cap)
	return &s
}

func (s *Storage) AddItem(i api.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.cache[i.Key]; !ok {
		s.cache[i.Key] = i.Value
		return nil
	}
	return ErrItemAlreadyExists
}

func (s *Storage) RemoveItem(key api.ItemKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.cache[key]; ok {
		delete(s.cache, key)
		return nil
	}
	return ErrNoSuchItem
}

func (s *Storage) GetItem(key api.ItemKey) (api.ItemValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if val, ok := s.cache[key]; ok {
		return val, nil
	}
	return "", ErrNoSuchItem
}

func (s *Storage) GetAllItems() (items []api.Item) {
	s.mu.Lock()

	for k, v := range s.cache {
		items = append(items, api.Item{Key: k, Value: v})
	}
	s.mu.Unlock()
	return items
}
