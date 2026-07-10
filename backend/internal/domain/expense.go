package domain

import (
	"errors"
	"strings"
)

type Expense struct {
	ID           string             `json:"id"`
	Description  string             `json:"description"`
	Amount       float64            `json:"amount"`
	Participants []string           `json:"participants"`
	Payments     map[string]float64 `json:"payments"`
}

func (e Expense) Validate(allowedParticipants []string) error {
	if strings.TrimSpace(e.ID) == "" {
		return errors.New("expense id is required")
	}

	if strings.TrimSpace(e.Description) == "" {
		return errors.New("expense description is required")
	}

	if e.Amount <= 0 {
		return errors.New("expense amount must be greater than zero")
	}

	if len(e.Participants) == 0 {
		return errors.New("expense must include at least one participant")
	}

	allowed := make(map[string]struct{}, len(allowedParticipants))
	for _, participant := range allowedParticipants {
		participant = strings.TrimSpace(participant)
		if participant != "" {
			allowed[participant] = struct{}{}
		}
	}

	for _, participant := range e.Participants {
		participant = strings.TrimSpace(participant)
		if participant == "" {
			return errors.New("expense participant names must not be empty")
		}
		if _, exists := allowed[participant]; !exists {
			return errors.New("expense participant must be a valid event participant")
		}
	}

	for payer, amountPaid := range e.Payments {
		payer = strings.TrimSpace(payer)
		if payer == "" {
			return errors.New("payment payer must not be empty")
		}
		if amountPaid < 0 {
			return errors.New("payment amount must not be negative")
		}
		if _, exists := allowed[payer]; !exists {
			return errors.New("payment payer must be a valid event participant")
		}
	}

	return nil
}

func (e Expense) CalculateShares() map[string]float64 {
	shares := make(map[string]float64)
	if len(e.Participants) == 0 {
		return shares
	}

	perPerson := e.Amount / float64(len(e.Participants))
	for _, participant := range e.Participants {
		shares[participant] = perPerson
	}

	for payer, amountPaid := range e.Payments {
		shares[payer] -= amountPaid
	}

	return shares
}
