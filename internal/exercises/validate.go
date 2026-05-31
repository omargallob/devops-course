package exercises

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationResult holds the outcome of validating user code output.
type ValidationResult struct {
	Passed         bool
	ActualOutput   string
	ExpectedOutput string
	Diff           string
	CompileError   string
}

// Validate compares actual output against an exercise's expected output
// using the exercise's configured validation mode.
func Validate(exercise *Exercise, actualOutput, compileError string) ValidationResult {
	if compileError != "" {
		return ValidationResult{
			Passed:       false,
			ActualOutput: actualOutput,
			CompileError: compileError,
		}
	}

	switch exercise.ValidationMode {
	case ValidationModeRegex:
		return validateRegex(exercise, actualOutput)
	default:
		return validateExact(exercise, actualOutput)
	}
}

func validateExact(exercise *Exercise, actualOutput string) ValidationResult {
	expected := exercise.ExpectedOutput
	if actualOutput == expected {
		return ValidationResult{
			Passed:         true,
			ActualOutput:   actualOutput,
			ExpectedOutput: expected,
		}
	}

	return ValidationResult{
		Passed:         false,
		ActualOutput:   actualOutput,
		ExpectedOutput: expected,
		Diff:           computeDiff(expected, actualOutput),
	}
}

func validateRegex(exercise *Exercise, actualOutput string) ValidationResult {
	pattern := exercise.ExpectedOutput
	matched, err := regexp.MatchString(pattern, actualOutput)
	if err != nil {
		return ValidationResult{
			Passed:         false,
			ActualOutput:   actualOutput,
			ExpectedOutput: pattern,
			Diff:           fmt.Sprintf("invalid regex pattern: %v", err),
		}
	}

	if matched {
		return ValidationResult{
			Passed:         true,
			ActualOutput:   actualOutput,
			ExpectedOutput: pattern,
		}
	}

	return ValidationResult{
		Passed:         false,
		ActualOutput:   actualOutput,
		ExpectedOutput: pattern,
		Diff:           fmt.Sprintf("output did not match pattern: %s", pattern),
	}
}

// computeDiff produces a simple line-by-line diff between expected and actual.
func computeDiff(expected, actual string) string {
	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(actual, "\n")

	var b strings.Builder
	maxLen := len(expectedLines)
	if len(actualLines) > maxLen {
		maxLen = len(actualLines)
	}

	for i := 0; i < maxLen; i++ {
		var exp, act string
		if i < len(expectedLines) {
			exp = expectedLines[i]
		}
		if i < len(actualLines) {
			act = actualLines[i]
		}

		if exp == act {
			fmt.Fprintf(&b, "  %s\n", exp)
		} else {
			if i < len(expectedLines) {
				fmt.Fprintf(&b, "- %s\n", exp)
			}
			if i < len(actualLines) {
				fmt.Fprintf(&b, "+ %s\n", act)
			}
		}
	}

	return b.String()
}
