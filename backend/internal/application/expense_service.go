package application

import (
	"expense-splitter/internal/domain"
)

type ExpenseUseCase struct {
	repository EventRepository
}

func NewExpenseUseCase(repository EventRepository) *ExpenseUseCase {
	return &ExpenseUseCase{repository: repository}
}

func (u *ExpenseUseCase) CreateEvent(event *domain.Event) error {
	return u.repository.CreateEvent(event)
}

func (u *ExpenseUseCase) GetEvent(id string) (*domain.Event, error) {
	return u.repository.GetEvent(id)
}

func (u *ExpenseUseCase) ListEvents() ([]*domain.Event, error) {
	return u.repository.ListEvents()
}

func (u *ExpenseUseCase) UpdateEvent(event *domain.Event) error {
	return u.repository.UpdateEvent(event)
}

func (u *ExpenseUseCase) DeleteEvent(id string) error {
	return u.repository.DeleteEvent(id)
}

func (u *ExpenseUseCase) CalculateBalances(eventID string) (map[string]float64, error) {
	event, err := u.repository.GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	return event.CalculateBalances(), nil
}
