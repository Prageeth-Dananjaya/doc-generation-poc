package domain

import "testing"

func TestCalculateSharesWithOnePayer(t *testing.T) {
	expense := Expense{
		Description:  "Dinner",
		Amount:       120,
		Participants: []string{"alice", "bob", "carol"},
		Payments: map[string]float64{
			"alice": 120,
		},
	}

	shares := expense.CalculateShares()

	if got := shares["alice"]; got != -80 {
		t.Fatalf("expected alice share to be %.2f, got %.2f", -80.0, got)
	}
	if got := shares["bob"]; got != 40 {
		t.Fatalf("expected bob share to be %.2f, got %.2f", 40.0, got)
	}
	if got := shares["carol"]; got != 40 {
		t.Fatalf("expected carol share to be %.2f, got %.2f", 40.0, got)
	}
}

func TestCalculateSharesWithoutPayer(t *testing.T) {
	expense := Expense{
		Description:  "Lunch",
		Amount:       100,
		Participants: []string{"dave", "erin"},
	}

	shares := expense.CalculateShares()

	if got := shares["dave"]; got != 50 {
		t.Fatalf("expected dave share to be %.2f, got %.2f", 50.0, got)
	}
	if got := shares["erin"]; got != 50 {
		t.Fatalf("expected erin share to be %.2f, got %.2f", 50.0, got)
	}
}

func TestExpenseValidateSucceeds(t *testing.T) {
	expense := Expense{
		ID:           "expense-1",
		Description:  "Dinner",
		Amount:       120,
		Participants: []string{"alice", "bob", "carol"},
		Payments: map[string]float64{
			"alice": 120,
		},
	}

	if err := expense.Validate([]string{"alice", "bob", "carol"}); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestExpenseValidateRejectsNegativeAmount(t *testing.T) {
	expense := Expense{
		ID:           "expense-1",
		Description:  "Dinner",
		Amount:       -5,
		Participants: []string{"alice", "bob"},
	}

	if err := expense.Validate([]string{"alice", "bob"}); err == nil {
		t.Fatal("expected validation error for negative amount")
	}
}

func TestExpenseValidateRejectsInvalidPayer(t *testing.T) {
	expense := Expense{
		ID:           "expense-1",
		Description:  "Dinner",
		Amount:       100,
		Participants: []string{"alice", "bob"},
		Payments: map[string]float64{
			"carol": 100,
		},
	}

	if err := expense.Validate([]string{"alice", "bob"}); err == nil {
		t.Fatal("expected validation error for payment payer outside participants")
	}
}
