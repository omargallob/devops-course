package exercises

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func strReader(s string) io.Reader {
	return strings.NewReader(s)
}

func TestGetExercise_Success(t *testing.T) {
	handler := &Handler{
		store: NewStore(),
	}

	r := chi.NewRouter()
	r.Get("/api/exercises/{exerciseId}", handler.GetExercise)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises/m01-hello-world", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp ExerciseAPIResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != "m01-hello-world" {
		t.Errorf("expected id 'm01-hello-world', got %q", resp.ID)
	}
	if resp.Title != "Hello World" {
		t.Errorf("expected title 'Hello World', got %q", resp.Title)
	}
	if resp.StarterCode == "" {
		t.Error("expected starter code to be non-empty")
	}
}

func TestGetExercise_NotFound(t *testing.T) {
	handler := &Handler{
		store: NewStore(),
	}

	r := chi.NewRouter()
	r.Get("/api/exercises/{exerciseId}", handler.GetExercise)

	req := httptest.NewRequest(http.MethodGet, "/api/exercises/nonexistent", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestValidateExercise_BadRequest_EmptyBody(t *testing.T) {
	handler := &Handler{
		store: NewStore(),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/validate", http.NoBody)
	w := httptest.NewRecorder()
	handler.ValidateExercise(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestValidateExercise_BadRequest_MissingExerciseID(t *testing.T) {
	handler := &Handler{
		store: NewStore(),
	}

	body := `{"code": "package main"}`
	req := httptest.NewRequest(http.MethodPost, "/api/validate", strReader(body))
	w := httptest.NewRecorder()
	handler.ValidateExercise(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestValidateExercise_NotFound(t *testing.T) {
	handler := &Handler{
		store: NewStore(),
	}

	body := `{"exerciseId": "nonexistent", "code": "package main"}`
	req := httptest.NewRequest(http.MethodPost, "/api/validate", strReader(body))
	w := httptest.NewRecorder()
	handler.ValidateExercise(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}
