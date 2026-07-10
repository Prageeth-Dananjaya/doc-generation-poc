package application

import (
	"errors"
	"sync"

	"expense-splitter/internal/domain"
)

type InMemoryEventRepository struct {
	mu     sync.RWMutex
	events map[string]*domain.Event
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		events: make(map[string]*domain.Event),
	}
}

func (r *InMemoryEventRepository) CreateEvent(event *domain.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[event.ID]; exists {
		return errors.New("event already exists")
	}

	r.events[event.ID] = event
	return nil
}

func (r *InMemoryEventRepository) GetEvent(id string) (*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return nil, ErrEventNotFound
	}

	return event, nil
}

func (r *InMemoryEventRepository) ListEvents() ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Event, 0, len(r.events))
	for _, event := range r.events {
		result = append(result, event)
	}

	return result, nil
}

func (r *InMemoryEventRepository) UpdateEvent(event *domain.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[event.ID]; !exists {
		return ErrEventNotFound
	}

	r.events[event.ID] = event
	return nil
}

func (r *InMemoryEventRepository) DeleteEvent(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[id]; !exists {
		return ErrEventNotFound
	}

	delete(r.events, id)
	return nil
}
