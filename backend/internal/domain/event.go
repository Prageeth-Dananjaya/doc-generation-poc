package domain

type Event struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Participants []string  `json:"participants"`
	Expenses     []Expense `json:"expenses"`
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
