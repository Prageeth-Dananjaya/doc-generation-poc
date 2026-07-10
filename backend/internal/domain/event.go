package domain

import (
	"errors"
	"strings"
)

type Event struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Participants []string  `json:"participants"`
	Expenses     []Expense `json:"expenses"`
}

func (e Event) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return errors.New("event id is required")
	}

	if strings.TrimSpace(e.Name) == "" {
		return errors.New("event name is required")
	}

	if len(e.Participants) == 0 {
		return errors.New("event must have at least one participant")
	}

	known := make(map[string]struct{}, len(e.Participants))
	for _, participant := range e.Participants {
		participant = strings.TrimSpace(participant)
		if participant == "" {
			return errors.New("participant names must not be empty")
		}
		if _, exists := known[participant]; exists {
			return errors.New("participant names must be unique")
		}
		known[participant] = struct{}{}
	}

	for _, expense := range e.Expenses {
		if err := expense.Validate(e.Participants); err != nil {
			return err
		}
	}

	return nil
}

func (e Event) HasParticipant(name string) bool {
	_, exists := e.participantSet()[strings.TrimSpace(name)]
	return exists
}

func (e Event) participantSet() map[string]struct{} {
	set := make(map[string]struct{}, len(e.Participants))
	for _, participant := range e.Participants {
		participant = strings.TrimSpace(participant)
		if participant != "" {
			set[participant] = struct{}{}
		}
	}
	return set
}

func (e Event) CalculateBalances() map[string]float64 {
	balances := make(map[string]float64, len(e.Participants))
	for _, participant := range e.Participants {
		balances[participant] = 0
	}

	for _, expense := range e.Expenses {
		if len(expense.Participants) == 0 {
			continue
		}

		perPersonShare := expense.Amount / float64(len(expense.Participants))
		for _, participant := range expense.Participants {
			balances[participant] += perPersonShare
		}

		for payer, amountPaid := range expense.Payments {
			if _, exists := balances[payer]; !exists {
				balances[payer] = 0
			}
			balances[payer] -= amountPaid
		}
	}

	return balances
}
