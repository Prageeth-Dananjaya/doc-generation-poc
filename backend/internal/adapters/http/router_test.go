package http

import (
	"bytes"
	"encoding/json"
	"expense-splitter/internal/application"
	"expense-splitter/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestRouter() http.Handler {
	repository := application.NewInMemoryEventRepository()
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
	_ = repository.CreateEvent(event)
	useCase := application.NewExpenseUseCase(repository)
	return NewRouter(useCase)
}

func TestGetBalancesEndpoint(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/events/event-1/balances", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var body map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	balances, ok := body["balances"]
	if !ok {
		t.Fatal("expected balances field")
	}

	if got := balances["alice"]; got != -80 {
		t.Fatalf("expected alice balance -80, got %.2f", got)
	}
}

func TestPostEventRejectsInvalidPayload(t *testing.T) {
	router := newTestRouter()
	invalidEvent := map[string]any{
		"name":         "Trip without ID",
		"participants": []string{"alice"},
	}
	payload, _ := json.Marshal(invalidEvent)

	req := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["error"] == "" {
		t.Fatal("expected error message")
	}
}

func TestPostEventCreatesEvent(t *testing.T) {
	router := newTestRouter()
	newEvent := &domain.Event{
		ID:           "event-2",
		Name:         "Weekend Trip",
		Participants: []string{"alice", "bob"},
		Expenses: []domain.Expense{{
			ID:           "expense-2",
			Description:  "Coffee",
			Amount:       20,
			Participants: []string{"alice", "bob"},
			Payments: map[string]float64{
				"alice": 20,
			},
		}},
	}
	payload, _ := json.Marshal(newEvent)

	req := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.Code)
	}

	var body domain.Event
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.ID != "event-2" {
		t.Fatalf("expected event id %s, got %s", "event-2", body.ID)
	}
}
