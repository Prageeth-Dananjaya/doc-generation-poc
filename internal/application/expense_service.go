package application

import "expense-splitter/internal/domain"

type ExpenseUseCase struct {
	event *domain.Event
}

func NewExpenseUseCase(event *domain.Event) *ExpenseUseCase {
	return &ExpenseUseCase{event: event}
}

func (u *ExpenseUseCase) CalculateBalances() map[string]float64 {
	return u.event.CalculateBalances()
}
