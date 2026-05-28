// Package playground proxies code execution requests to the Go Playground
// (https://go.dev/_/compile). This allows the frontend to run user-submitted
// Go code without hosting a custom sandbox. The playground URL is exported on
// CompileHandler for test mocking.
package playground

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const goPlaygroundURL = "https://go.dev/_/compile"

// CompileRequest is the request body for the compile endpoint.
type CompileRequest struct {
	Body    string `json:"body"`
	Version int    `json:"version,omitempty"`
}

// CompileResponse is the response from the Go Playground.
type CompileResponse struct {
	Errors string `json:"Errors"`
	Events []struct {
		Message string `json:"Message"`
		Kind    string `json:"Kind"`
		Delay   int    `json:"Delay"`
	} `json:"Events"`
	Status      int    `json:"Status"`
	IsTest      bool   `json:"IsTest"`
	TestsFailed int    `json:"TestsFailed"`
	VetOK       bool   `json:"VetOK"`
	VetErrors   string `json:"VetErrors"`
}

// CompileHandler handles compile requests by proxying to the Go Playground.
type CompileHandler struct {
	logger        *slog.Logger
	client        *http.Client
	playgroundURL string
}

// NewCompileHandler creates a new CompileHandler.
func NewCompileHandler(logger *slog.Logger) *CompileHandler {
	return &CompileHandler{
		logger:        logger,
		playgroundURL: goPlaygroundURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Handle proxies a compile request to the Go Playground.
func (h *CompileHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req CompileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Body) == "" {
		http.Error(w, `{"error":"empty code body"}`, http.StatusBadRequest)
		return
	}

	// Proxy to Go Playground.
	// NOTE: req.Body is passed as a raw form value. The Go Playground accepts
	// un-encoded source code in the "body" field, so we do not URL-encode it.
	form := strings.NewReader("version=2&body=" + req.Body)
	proxyReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, h.playgroundURL, form)
	if err != nil {
		h.logger.Error("failed to create proxy request", "error", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	proxyReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := h.client.Do(proxyReq)
	if err != nil {
		h.logger.Error("playground request failed", "error", err)
		http.Error(w, `{"error":"playground unavailable"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
