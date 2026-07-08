package main

import (
	httpadapter "expense-splitter/internal/adapters/http"
	"expense-splitter/internal/application"
	"expense-splitter/internal/domain"
	"log"
	nethttp "net/http"
)

func main() {
	event := &domain.Event{
		ID:           "event-1",
		Name:         "Sample Trip",
		Participants: []string{"alice", "bob", "carol"},
		Expenses: []domain.Expense{{
			Description:  "Dinner",
			Amount:       120,
			Participants: []string{"alice", "bob", "carol"},
			Payments: map[string]float64{
				"alice": 120,
			},
		}},
	}

	useCase := application.NewExpenseUseCase(event)
	router := httpadapter.NewRouter(useCase)

	port := ":8081"
	log.Printf("server listening on %s", port)
	if err := nethttp.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
