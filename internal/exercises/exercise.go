// Package exercises provides exercise definitions, storage, and validation
// for the interactive Go course.
package exercises

// ValidationMode determines how exercise output is compared.
type ValidationMode string

// Validation modes for exercise output comparison.
const (
	ValidationModeExact ValidationMode = "exact"
	ValidationModeRegex ValidationMode = "regex"
)

// Exercise defines a single coding exercise.
type Exercise struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Instructions   string         `json:"instructions"`
	StarterCode    string         `json:"starterCode"`
	Hint           string         `json:"hint,omitempty"`
	ExpectedOutput string         `json:"expectedOutput"`
	ValidationMode ValidationMode `json:"validationMode"`
}
