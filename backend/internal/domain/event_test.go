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
