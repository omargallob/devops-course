package exercises

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/omargallob/devops-course/internal/playground"
)

// Handler serves exercise-related HTTP endpoints.
type Handler struct {
	logger  *slog.Logger
	store   *Store
	compile *playground.CompileHandler
}

// NewHandler creates a Handler with its dependencies.
func NewHandler(logger *slog.Logger, compile *playground.CompileHandler) *Handler {
	return &Handler{
		logger:  logger,
		store:   NewStore(),
		compile: compile,
	}
}

// GetExercise handles GET /api/exercises/{exerciseId}.
func (h *Handler) GetExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID := chi.URLParam(r, "exerciseId")
	if exerciseID == "" {
		writeError(w, http.StatusBadRequest, "missing exercise ID")
		return
	}

	exercise, err := h.store.Get(exerciseID)
	if err != nil {
		writeError(w, http.StatusNotFound, "exercise not found")
		return
	}

	// Return exercise without expected output (prevent cheating)
	resp := ExerciseAPIResponse{
		ID:             exercise.ID,
		Title:          exercise.Title,
		Instructions:   exercise.Instructions,
		StarterCode:    exercise.StarterCode,
		Hint:           exercise.Hint,
		ValidationMode: string(exercise.ValidationMode),
	}

	writeJSON(w, http.StatusOK, resp)
}

// ValidateExercise handles POST /api/validate.
func (h *Handler) ValidateExercise(w http.ResponseWriter, r *http.Request) {
	var req ValidateAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Code) == "" {
		writeError(w, http.StatusBadRequest, "empty code body")
		return
	}
	if req.ExerciseID == "" {
		writeError(w, http.StatusBadRequest, "missing exercise ID")
		return
	}

	exercise, err := h.store.Get(req.ExerciseID)
	if err != nil {
		writeError(w, http.StatusNotFound, "exercise not found")
		return
	}

	// Compile and run user code via the playground proxy
	compileResp, err := h.compile.CompileCode(r.Context(), req.Code)
	if err != nil {
		h.logger.Error("compile failed", "error", err, "exerciseId", req.ExerciseID)
		writeError(w, http.StatusBadGateway, "playground unavailable")
		return
	}

	// Extract output or compile error
	var actualOutput, compileError string
	if compileResp.Errors != "" {
		compileError = compileResp.Errors
	} else {
		actualOutput = extractOutput(compileResp)
	}

	// Validate
	result := Validate(&exercise, actualOutput, compileError)

	resp := ValidateAPIResponse{
		Passed:     result.Passed,
		ExerciseID: req.ExerciseID,
	}
	if result.ActualOutput != "" {
		resp.ActualOutput = &result.ActualOutput
	}
	if !result.Passed {
		resp.ExpectedOutput = &result.ExpectedOutput
		if result.Diff != "" {
			resp.Diff = &result.Diff
		}
		if result.CompileError != "" {
			resp.CompileError = &result.CompileError
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

// extractOutput concatenates all stdout events from a compile response.
func extractOutput(resp *playground.CompileResponse) string {
	var b strings.Builder
	for _, ev := range resp.Events {
		if ev.Kind == "stdout" || ev.Kind == "" {
			b.WriteString(ev.Message)
		}
	}
	return b.String()
}

// --- API request/response types ---

// ExerciseAPIResponse is the JSON response for GET /api/exercises/{id}.
type ExerciseAPIResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Instructions   string `json:"instructions"`
	StarterCode    string `json:"starterCode"`
	Hint           string `json:"hint,omitempty"`
	ValidationMode string `json:"validationMode"`
}

// ValidateAPIRequest is the JSON request for POST /api/validate.
type ValidateAPIRequest struct {
	ExerciseID string `json:"exerciseId"`
	Code       string `json:"code"`
}

// ValidateAPIResponse is the JSON response for POST /api/validate.
type ValidateAPIResponse struct {
	Passed         bool    `json:"passed"`
	ExerciseID     string  `json:"exerciseId"`
	ActualOutput   *string `json:"actualOutput,omitempty"`
	ExpectedOutput *string `json:"expectedOutput,omitempty"`
	Diff           *string `json:"diff,omitempty"`
	CompileError   *string `json:"compileError,omitempty"`
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
