package template

import (
	"os"
	"testing"
)

func TestRenderSimple(t *testing.T) {
	eng := NewEngine(nil)
	eng.Merge(map[string]string{
		"name": "John",
		"age":  "30",
	})

	result := eng.Render("Hello ${name}, you are ${age} years old")
	expected := "Hello John, you are 30 years old"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderMissing(t *testing.T) {
	eng := NewEngine(nil)

	result := eng.Render("Hello ${missing}")
	expected := "Hello ${missing}"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderNoVariables(t *testing.T) {
	eng := NewEngine(nil)

	result := eng.Render("Hello World")
	expected := "Hello World"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderMap(t *testing.T) {
	eng := NewEngine(nil)
	eng.Merge(map[string]string{
		"host": "localhost",
		"port": "8080",
	})

	input := map[string]string{
		"url": "http://${host}:${port}/api",
	}
	result := eng.RenderMap(input)

	if result["url"] != "http://localhost:8080/api" {
		t.Errorf("expected http://localhost:8080/api, got %s", result["url"])
	}
}

func TestLoadEnv(t *testing.T) {
	os.Setenv("HTTPFORGE_TEST_VAR", "test_value")
	defer os.Unsetenv("HTTPFORGE_TEST_VAR")

	eng := NewEngine(nil)
	eng.LoadEnv("HTTPFORGE_")

	vars := eng.Vars()
	if vars["TEST_VAR"] != "test_value" {
		t.Errorf("expected test_value, got %s", vars["TEST_VAR"])
	}
}

func TestValidate(t *testing.T) {
	eng := NewEngine(map[string]string{
		"name": "John",
	})

	// All variables available
	missing := eng.Validate("Hello ${name}")
	if len(missing) != 0 {
		t.Errorf("expected no missing variables, got %v", missing)
	}

	// Missing variable
	missing = eng.Validate("Hello ${name} from ${city}")
	if len(missing) != 1 || missing[0] != "city" {
		t.Errorf("expected [city], got %v", missing)
	}
}

func TestHasVariables(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Hello ${name}", true},
		{"Hello World", false},
		{"${a}${b}", true},
		{"", false},
	}

	for _, tt := range tests {
		result := HasVariables(tt.input)
		if result != tt.expected {
			t.Errorf("HasVariables(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestListVariables(t *testing.T) {
	vars := ListVariables("Hello ${name}, you are ${age} from ${city}")
	if len(vars) != 3 {
		t.Errorf("expected 3 variables, got %d", len(vars))
	}
}

func TestExtractVarName(t *testing.T) {
	name := ExtractVarName("${my_var}")
	if name != "my_var" {
		t.Errorf("expected my_var, got %s", name)
	}
}
