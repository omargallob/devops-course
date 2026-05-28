package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseWriterCapturesStatus(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{"200 OK", http.StatusOK},
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
		{"201 Created", http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			rw.WriteHeader(tt.code)

			if rw.statusCode != tt.code {
				t.Errorf("expected status %d, got %d", tt.code, rw.statusCode)
			}
		})
	}
}

func TestResponseWriterDefaultStatus(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	// Without calling WriteHeader, default should be 200
	if rw.statusCode != http.StatusOK {
		t.Errorf("expected default status 200, got %d", rw.statusCode)
	}
}

func TestSlogMiddlewareDoesNotPanic(t *testing.T) {
	logger := testLogger()
	mw := slogMiddleware(logger)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()

	// Should not panic
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
