package http

import (
	"encoding/json"
	"errors"
	"expense-splitter/internal/application"
	"expense-splitter/internal/domain"
	"net/http"
	"strings"
)

func NewRouter(useCase *application.ExpenseUseCase) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/balances", func(w http.ResponseWriter, r *http.Request) {
		balances, err := useCase.CalculateBalances("event-1")
		if err != nil {
			writeJSONError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"balances": balances})
	})

	mux.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			events, err := useCase.ListEvents()
			if err != nil {
				writeJSONError(w, http.StatusInternalServerError, err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"events": events})
		case http.MethodPost:
			event := &domain.Event{}
			if err := json.NewDecoder(r.Body).Decode(event); err != nil {
				writeJSONError(w, http.StatusBadRequest, err)
				return
			}
			if err := useCase.CreateEvent(event); err != nil {
				writeJSONError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusCreated, event)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/events/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/events/")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			event, err := useCase.GetEvent(id)
			if err != nil {
				writeRepositoryError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, event)
		case http.MethodPut:
			event := &domain.Event{}
			if err := json.NewDecoder(r.Body).Decode(event); err != nil {
				writeJSONError(w, http.StatusBadRequest, err)
				return
			}
			if event.ID == "" {
				event.ID = id
			}
			if event.ID != id {
				writeJSONError(w, http.StatusBadRequest, errors.New("event ID does not match URL"))
				return
			}
			if err := useCase.UpdateEvent(event); err != nil {
				writeRepositoryError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, event)
		case http.MethodDelete:
			if err := useCase.DeleteEvent(id); err != nil {
				writeRepositoryError(w, err)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeJSONError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeRepositoryError(w http.ResponseWriter, err error) {
	if errors.Is(err, application.ErrEventNotFound) {
		writeJSONError(w, http.StatusNotFound, err)
		return
	}
	writeJSONError(w, http.StatusBadRequest, err)
}
