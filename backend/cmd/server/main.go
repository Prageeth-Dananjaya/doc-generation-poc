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
			ID:           "expense-1",
			Description:  "Dinner",
			Amount:       120,
			Participants: []string{"alice", "bob", "carol"},
			Payments: map[string]float64{
				"alice": 120,
			},
		}},
	}

	repository := application.NewInMemoryEventRepository()
	if err := repository.CreateEvent(event); err != nil {
		log.Fatalf("failed to seed event: %v", err)
	}

	useCase := application.NewExpenseUseCase(repository)
	router := httpadapter.NewRouter(useCase)

	port := ":8081"
	log.Printf("server listening on %s", port)
	if err := nethttp.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
