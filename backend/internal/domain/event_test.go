package domain

import "testing"

func TestCalculateBalances(t *testing.T) {
	event := Event{
		Name:         "Trip",
		Participants: []string{"alice", "bob", "carol"},
		Expenses: []Expense{
			{
				Description:  "Dinner",
				Amount:       120,
				Participants: []string{"alice", "bob", "carol"},
				Payments: map[string]float64{
					"alice": 120,
				},
			},
		},
	}

	balances := event.CalculateBalances()

	if got := balances["alice"]; got != -80 {
		t.Fatalf("expected alice balance to be %.2f, got %.2f", -80.0, got)
	}
	if got := balances["bob"]; got != 40 {
		t.Fatalf("expected bob balance to be %.2f, got %.2f", 40.0, got)
	}
	if got := balances["carol"]; got != 40 {
		t.Fatalf("expected carol balance to be %.2f, got %.2f", 40.0, got)
	}
}

func TestEventValidateSucceeds(t *testing.T) {
	event := Event{
		ID:           "event-1",
		Name:         "Trip",
		Participants: []string{"alice", "bob", "carol"},
		Expenses: []Expense{
			{
				ID:           "expense-1",
				Description:  "Dinner",
				Amount:       120,
				Participants: []string{"alice", "bob", "carol"},
				Payments: map[string]float64{
					"alice": 120,
				},
			},
		},
	}

	if err := event.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestEventValidateRejectsDuplicateParticipants(t *testing.T) {
	event := Event{
		ID:           "event-1",
		Name:         "Trip",
		Participants: []string{"alice", "alice"},
	}

	if err := event.Validate(); err == nil {
		t.Fatal("expected validation error for duplicate participants")
	}
}

func TestEventValidateRejectsInvalidExpenseParticipant(t *testing.T) {
	event := Event{
		ID:           "event-1",
		Name:         "Trip",
		Participants: []string{"alice", "bob"},
		Expenses: []Expense{
			{
				ID:           "expense-1",
				Description:  "Dinner",
				Amount:       100,
				Participants: []string{"alice", "carol"},
				Payments: map[string]float64{
					"alice": 100,
				},
			},
		},
	}

	if err := event.Validate(); err == nil {
		t.Fatal("expected validation error for expense participant outside event")
	}
}
