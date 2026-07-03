package assert

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/httpforge/models"
)

func TestValidateStatusCode(t *testing.T) {
	resp := &models.Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"status": "ok"}`,
	}

	assertions := []models.Assertion{
		{Field: "status_code", Op: "eq", Value: "200"},
		{Field: "status", Op: "eq", Value: "200"},
	}

	results := Validate(resp, assertions)
	for _, r := range results {
		if !r.Passed {
			t.Errorf("assertion failed: %s", r.Message)
		}
	}
}

func TestValidateHeader(t *testing.T) {
	resp := &models.Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}

	assertions := []models.Assertion{
		{Field: "header", Op: "eq", Value: "Content-Type"},
	}

	// This tests that we can access headers
	results := Validate(resp, assertions)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestValidateJSONPath(t *testing.T) {
	resp := &models.Response{
		StatusCode: 200,
		Body:       `{"user": {"name": "John", "age": 30}}`,
	}

	assertions := []models.Assertion{
		{Field: "json_path", Op: "eq", Value: "user.name"},
		{Field: "json_path", Op: "eq", Value: "user.age"},
	}

	results := Validate(resp, assertions)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// user.name should be "John"
	if results[0].Passed {
		t.Error("expected user.name assertion to fail (empty value)")
	}
}

func TestEvaluateOpContains(t *testing.T) {
	if !evaluateOp("hello world", "contains", "world") {
		t.Error("expected contains to return true")
	}
	if evaluateOp("hello world", "contains", "xyz") {
		t.Error("expected contains to return false")
	}
}

func TestEvaluateOpNumeric(t *testing.T) {
	if !evaluateOp("100", "gt", "50") {
		t.Error("expected gt to return true")
	}
	if evaluateOp("100", "lt", "50") {
		t.Error("expected lt to return false")
	}
	if !evaluateOp("100", "gte", "100") {
		t.Error("expected gte to return true")
	}
	if !evaluateOp("100", "lte", "100") {
		t.Error("expected lte to return true")
	}
}

func TestEvaluateOpNegative(t *testing.T) {
	if !evaluateOp("hello", "ne", "world") {
		t.Error("expected ne to return true")
	}
	if evaluateOp("hello", "ne", "hello") {
		t.Error("expected ne to return false")
	}
}

func TestAllPassed(t *testing.T) {
	results := []models.AssertionResult{
		{Passed: true},
		{Passed: true},
	}
	if !AllPassed(results) {
		t.Error("expected AllPassed to return true")
	}

	results = append(results, models.AssertionResult{Passed: false})
	if AllPassed(results) {
		t.Error("expected AllPassed to return false")
	}
}

func TestFailedCount(t *testing.T) {
	results := []models.AssertionResult{
		{Passed: true},
		{Passed: false},
		{Passed: false},
	}
	if FailedCount(results) != 2 {
		t.Errorf("expected 2 failures, got %d", FailedCount(results))
	}
}
