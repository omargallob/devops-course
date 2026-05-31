package exercises

import (
	"strings"
	"testing"
)

func TestValidate_ExactMatch_Pass(t *testing.T) {
	ex := Exercise{
		ID:             "test-exact",
		ExpectedOutput: "Hello, World!\n",
		ValidationMode: ValidationModeExact,
	}

	result := Validate(&ex, "Hello, World!\n", "")
	if !result.Passed {
		t.Errorf("expected pass, got fail: diff=%s", result.Diff)
	}
}

func TestValidate_ExactMatch_Fail(t *testing.T) {
	ex := Exercise{
		ID:             "test-exact",
		ExpectedOutput: "Hello, World!\n",
		ValidationMode: ValidationModeExact,
	}

	result := Validate(&ex, "hello world\n", "")
	if result.Passed {
		t.Error("expected fail, got pass")
	}
	if result.Diff == "" {
		t.Error("expected diff to be non-empty")
	}
	if result.ExpectedOutput != "Hello, World!\n" {
		t.Errorf("expected output mismatch: got %q", result.ExpectedOutput)
	}
}

func TestValidate_CompileError(t *testing.T) {
	ex := Exercise{
		ID:             "test-compile",
		ExpectedOutput: "anything",
		ValidationMode: ValidationModeExact,
	}

	result := Validate(&ex, "", "prog.go:5:1: expected declaration, found '}'")
	if result.Passed {
		t.Error("expected fail on compile error")
	}
	if result.CompileError == "" {
		t.Error("expected compile error to be set")
	}
}

func TestValidate_Regex_Pass(t *testing.T) {
	ex := Exercise{
		ID:             "test-regex",
		ExpectedOutput: `^Hello, \w+!$`,
		ValidationMode: ValidationModeRegex,
	}

	result := Validate(&ex, "Hello, Gopher!", "")
	if !result.Passed {
		t.Errorf("expected pass, got fail: diff=%s", result.Diff)
	}
}

func TestValidate_Regex_Fail(t *testing.T) {
	ex := Exercise{
		ID:             "test-regex",
		ExpectedOutput: `^Hello, \w+!$`,
		ValidationMode: ValidationModeRegex,
	}

	result := Validate(&ex, "Goodbye!", "")
	if result.Passed {
		t.Error("expected fail, got pass")
	}
}

func TestValidate_Regex_InvalidPattern(t *testing.T) {
	ex := Exercise{
		ID:             "test-regex-bad",
		ExpectedOutput: `[invalid`,
		ValidationMode: ValidationModeRegex,
	}

	result := Validate(&ex, "anything", "")
	if result.Passed {
		t.Error("expected fail on invalid regex")
	}
	if result.Diff == "" {
		t.Error("expected diff to contain error message")
	}
}

func TestStore_Get(t *testing.T) {
	store := NewStore()

	ex, err := store.Get("m01-hello-world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex.Title != "Hello World" {
		t.Errorf("expected title 'Hello World', got %q", ex.Title)
	}
}

func TestStore_Get_NotFound(t *testing.T) {
	store := NewStore()

	_, err := store.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent exercise")
	}
}

func TestStore_List(t *testing.T) {
	store := NewStore()

	exercises := store.List()
	if len(exercises) == 0 {
		t.Error("expected at least one exercise")
	}
}

func TestComputeDiff(t *testing.T) {
	diff := computeDiff("line1\nline2\n", "line1\nchanged\n")
	if diff == "" {
		t.Error("expected non-empty diff")
	}
	// Should contain the differing line
	if !strings.Contains(diff, "+ changed") {
		t.Errorf("diff should contain '+ changed', got: %s", diff)
	}
	if !strings.Contains(diff, "- line2") {
		t.Errorf("diff should contain '- line2', got: %s", diff)
	}
}
