package playground

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestCompileHandler_EmptyBody(t *testing.T) {
	handler := NewCompileHandler(testLogger())

	body, _ := json.Marshal(CompileRequest{Body: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/compile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "empty code body") {
		t.Errorf("expected error about empty body, got %q", w.Body.String())
	}
}

func TestCompileHandler_WhitespaceBody(t *testing.T) {
	handler := NewCompileHandler(testLogger())

	body, _ := json.Marshal(CompileRequest{Body: "   \n\t  "})
	req := httptest.NewRequest(http.MethodPost, "/api/compile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for whitespace-only body, got %d", w.Code)
	}
}

func TestCompileHandler_InvalidJSON(t *testing.T) {
	handler := NewCompileHandler(testLogger())

	req := httptest.NewRequest(http.MethodPost, "/api/compile", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "invalid request body") {
		t.Errorf("expected error about invalid JSON, got %q", w.Body.String())
	}
}

func TestCompileHandler_ProxiesToPlayground(t *testing.T) {
	// Mock Go Playground server
	mockPlayground := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if ct != "application/x-www-form-urlencoded" {
			t.Errorf("expected form content type, got %q", ct)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CompileResponse{
			Events: []struct {
				Message string `json:"Message"`
				Kind    string `json:"Kind"`
				Delay   int    `json:"Delay"`
			}{
				{Message: "Hello, World!\n", Kind: "stdout", Delay: 0},
			},
			Status: 0,
			VetOK:  true,
		})
	}))
	defer mockPlayground.Close()

	handler := NewCompileHandler(testLogger())
	handler.playgroundURL = mockPlayground.URL

	body, _ := json.Marshal(CompileRequest{Body: `package main; func main() { fmt.Println("Hello") }`})
	req := httptest.NewRequest(http.MethodPost, "/api/compile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp CompileResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Events) == 0 {
		t.Fatal("expected events in response")
	}
	if resp.Events[0].Message != "Hello, World!\n" {
		t.Errorf("expected 'Hello, World!' event, got %q", resp.Events[0].Message)
	}
}

func TestCompileHandler_PlaygroundError(t *testing.T) {
	// Mock server that returns compile errors
	mockPlayground := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CompileResponse{
			Errors: "./prog.go:3:1: expected declaration, got '}'",
			Status: 2,
		})
	}))
	defer mockPlayground.Close()

	handler := NewCompileHandler(testLogger())
	handler.playgroundURL = mockPlayground.URL

	body, _ := json.Marshal(CompileRequest{Body: `package main; }`})
	req := httptest.NewRequest(http.MethodPost, "/api/compile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 (playground returns errors in body), got %d", w.Code)
	}

	var resp CompileResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Errors == "" {
		t.Error("expected compile errors in response")
	}
}

func TestCompileHandler_PlaygroundDown(t *testing.T) {
	handler := NewCompileHandler(testLogger())
	// Point to a non-existent server
	handler.playgroundURL = "http://127.0.0.1:1"

	body, _ := json.Marshal(CompileRequest{Body: `package main; func main() {}`})
	req := httptest.NewRequest(http.MethodPost, "/api/compile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("expected status 502 when playground is down, got %d", w.Code)
	}
}

func TestNewCompileHandler(t *testing.T) {
	logger := testLogger()
	handler := NewCompileHandler(logger)

	if handler.logger != logger {
		t.Error("expected logger to be set")
	}
	if handler.client == nil {
		t.Error("expected HTTP client to be set")
	}
	if handler.playgroundURL != goPlaygroundURL {
		t.Errorf("expected default playground URL %q, got %q", goPlaygroundURL, handler.playgroundURL)
	}
}
