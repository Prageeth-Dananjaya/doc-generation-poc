package http

import (
	"encoding/json"
	"expense-splitter/internal/application"
	"expense-splitter/internal/domain"
	"net/http"
)

func NewRouter(useCase *application.ExpenseUseCase) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/balances", func(w http.ResponseWriter, r *http.Request) {
		balances := useCase.CalculateBalances()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"balances": balances})
	})

	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		event := &domain.Event{}
		if err := json.NewDecoder(r.Body).Decode(event); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		useCase = application.NewExpenseUseCase(event)
		balances := useCase.CalculateBalances()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"balances": balances})
	})

	return mux
}
