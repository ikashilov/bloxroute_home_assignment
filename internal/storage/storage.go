package storage

import (
	"assignmentapp/internal/api"
	"errors"
	orderedmap "github.com/wk8/go-ordered-map"
	"sync"
)

var (
	ErrItemAlreadyExists = errors.New("item already exists")
	ErrNoSuchItem        = errors.New("no such item")
)

type Storage struct {
	mu    sync.RWMutex
	cache *orderedmap.OrderedMap
}

func New(cap int) *Storage {
	var s Storage
	s.cache = orderedmap.New()
	return &s
}

func (s *Storage) AddItem(i api.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.cache.Get(i.Key); !ok {
		s.cache.Set(i.Key, i.Value)
		return nil
	}
	return ErrItemAlreadyExists
}

func (s *Storage) RemoveItem(key api.ItemKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.cache.Get(key); ok {
		s.cache.Delete(key)
		return nil
	}
	return ErrNoSuchItem
}

func (s *Storage) GetItem(key api.ItemKey) (api.ItemValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if val, ok := s.cache.Get(key); ok {
		return val.(api.ItemValue), nil
	}
	return api.EmptyItemValue, ErrNoSuchItem
}

func (s *Storage) GetAllItems() (items []api.Item) {
	s.mu.Lock()

	for pair := s.cache.Oldest(); pair != nil; pair = pair.Next() {
		items = append(items, api.Item{
			Key: pair.Key.(api.ItemKey),
			Value: pair.Value.(api.ItemValue),
		})
	}
	s.mu.Unlock()
	return items
}
