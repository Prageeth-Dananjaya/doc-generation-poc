package domain

type Expense struct {
	ID           string             `json:"id"`
	Description  string             `json:"description"`
	Amount       float64            `json:"amount"`
	Participants []string           `json:"participants"`
	Payments     map[string]float64 `json:"payments"`
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
