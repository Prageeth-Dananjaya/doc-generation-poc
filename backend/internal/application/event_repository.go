package application

import (
	"errors"

	"expense-splitter/internal/domain"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type EventRepository interface {
	CreateEvent(event *domain.Event) error
	GetEvent(id string) (*domain.Event, error)
	ListEvents() ([]*domain.Event, error)
	UpdateEvent(event *domain.Event) error
	DeleteEvent(id string) error
}
