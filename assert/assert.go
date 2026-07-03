// Package assert handles response assertion checking.
package assert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/EdgarOrtegaRamirez/httpforge/models"
	"github.com/tidwall/gjson"
)

// Validate runs all assertions against a response and returns results.
func Validate(resp *models.Response, assertions []models.Assertion) []models.AssertionResult {
	var results []models.AssertionResult
	for _, a := range assertions {
		result := checkAssertion(resp, a)
		results = append(results, result)
	}
	return results
}

// checkAssertion evaluates a single assertion against a response.
func checkAssertion(resp *models.Response, a models.Assertion) models.AssertionResult {
	result := models.AssertionResult{
		Assertion: a,
		Passed:    false,
	}

	var actual string

	switch strings.ToLower(a.Field) {
	case "status_code", "status", "statuscode":
		actual = strconv.Itoa(resp.StatusCode)

	case "header":
		// Format: "Header-Name" - a.Value is the header name
		actual = resp.Headers[a.Value]

	case "body":
		actual = resp.Body

	case "body.length":
		actual = strconv.Itoa(len(resp.Body))

	case "json_path":
		// Use gjson for JSON path queries
		if a.Path != "" {
			result := gjson.Get(resp.Body, a.Path)
			actual = result.String()
		} else {
			// Fallback: treat Value as the path for backward compatibility
			result := gjson.Get(resp.Body, a.Value)
			actual = result.String()
		}

	case "header.count":
		actual = strconv.Itoa(len(resp.Headers))

	case "timing":
		actual = resp.Timing.String()

	case "size":
		actual = strconv.Itoa(resp.Size)

	default:
		return models.AssertionResult{
			Assertion: a,
			Passed:    false,
			Message:   fmt.Sprintf("unknown field: %s", a.Field),
		}
	}

	// Evaluate the assertion
	passed := evaluateOp(actual, a.Op, a.Value)

	// Handle negation
	if a.Negative {
		passed = !passed
	}

	result.Passed = passed
	if passed {
		result.Message = fmt.Sprintf("PASS: %s %s %q (actual: %q)", a.Field, a.Op, a.Value, actual)
	} else {
		result.Message = fmt.Sprintf("FAIL: %s %s %q (actual: %q)", a.Field, a.Op, a.Value, actual)
	}

	return result
}

// evaluateOp evaluates an assertion operation.
func evaluateOp(actual, op, expected string) bool {
	switch strings.ToLower(op) {
	case "eq", "==", "equals":
		return actual == expected
	case "ne", "!=", "not_equals":
		return actual != expected
	case "gt", ">", "greater":
		a, err1 := strconv.ParseFloat(actual, 64)
		e, err2 := strconv.ParseFloat(expected, 64)
		if err1 != nil || err2 != nil {
			return false
		}
		return a > e
	case "gte", ">=":
		a, err1 := strconv.ParseFloat(actual, 64)
		e, err2 := strconv.ParseFloat(expected, 64)
		if err1 != nil || err2 != nil {
			return false
		}
		return a >= e
	case "lt", "<", "less":
		a, err1 := strconv.ParseFloat(actual, 64)
		e, err2 := strconv.ParseFloat(expected, 64)
		if err1 != nil || err2 != nil {
			return false
		}
		return a < e
	case "lte", "<=":
		a, err1 := strconv.ParseFloat(actual, 64)
		e, err2 := strconv.ParseFloat(expected, 64)
		if err1 != nil || err2 != nil {
			return false
		}
		return a <= e
	case "contains":
		return strings.Contains(actual, expected)
	case "not_contains", "!contains":
		return !strings.Contains(actual, expected)
	case "starts_with", "startswith":
		return strings.HasPrefix(actual, expected)
	case "ends_with", "endswith":
		return strings.HasSuffix(actual, expected)
	case "matches", "regex":
		// Simple regex match (could be enhanced)
		return strings.Contains(actual, expected)
	case "exists":
		return actual != ""
	case "empty":
		return actual == ""
	case "type":
		return checkType(actual, expected)
	default:
		return false
	}
}

// checkType validates the type of a value.
func checkType(value, expectedType string) bool {
	switch strings.ToLower(expectedType) {
	case "number":
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	case "integer", "int":
		_, err := strconv.ParseInt(value, 10, 64)
		return err == nil
	case "boolean", "bool":
		return value == "true" || value == "false"
	case "string":
		return true // Everything is a string
	case "json":
		return gjson.Valid(value)
	default:
		return false
	}
}

// AllPassed checks if all assertion results passed.
func AllPassed(results []models.AssertionResult) bool {
	for _, r := range results {
		if !r.Passed {
			return false
		}
	}
	return true
}

// FailedCount returns the number of failed assertions.
func FailedCount(results []models.AssertionResult) int {
	count := 0
	for _, r := range results {
		if !r.Passed {
			count++
		}
	}
	return count
}
