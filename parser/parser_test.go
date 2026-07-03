package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRequestFileYAML(t *testing.T) {
	content := `
request:
  name: Test GET
  method: GET
  url: https://httpbin.org/get
  headers:
    Accept: application/json
  query_params:
    foo: bar
  assertions:
    - field: status_code
      op: eq
      value: "200"
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	config, err := LoadRequestFile(path)
	if err != nil {
		t.Fatalf("LoadRequestFile failed: %v", err)
	}

	if config.Request.Method != "GET" {
		t.Errorf("expected method GET, got %s", config.Request.Method)
	}
	if config.Request.URL != "https://httpbin.org/get" {
		t.Errorf("expected URL https://httpbin.org/get, got %s", config.Request.URL)
	}
	if config.Request.Headers["Accept"] != "application/json" {
		t.Errorf("expected Accept header, got %v", config.Request.Headers)
	}
	if len(config.Request.Assertions) != 1 {
		t.Errorf("expected 1 assertion, got %d", len(config.Request.Assertions))
	}
}

func TestLoadRequestFileJSON(t *testing.T) {
	content := `{
  "request": {
    "name": "Test POST",
    "method": "POST",
    "url": "https://httpbin.org/post",
    "body": "{\"key\": \"value\"}",
    "body_type": "json"
  }
}`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	config, err := LoadRequestFile(path)
	if err != nil {
		t.Fatalf("LoadRequestFile failed: %v", err)
	}

	if config.Request.Method != "POST" {
		t.Errorf("expected method POST, got %s", config.Request.Method)
	}
	if config.Request.BodyType != "json" {
		t.Errorf("expected body_type json, got %s", config.Request.BodyType)
	}
}

func TestLoadRequestFileDefaultMethod(t *testing.T) {
	content := `
request:
  url: https://httpbin.org/get
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	config, err := LoadRequestFile(path)
	if err != nil {
		t.Fatalf("LoadRequestFile failed: %v", err)
	}

	if config.Request.Method != "GET" {
		t.Errorf("expected default method GET, got %s", config.Request.Method)
	}
}

func TestLoadRequestFileNotFound(t *testing.T) {
	_, err := LoadRequestFile("/nonexistent/file.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoadRequestFileInvalidYAML(t *testing.T) {
	content := `invalid: [yaml: content`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadRequestFile(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseRequestURL(t *testing.T) {
	req, err := ParseRequest("https://httpbin.org/get")
	if err != nil {
		t.Fatalf("ParseRequest failed: %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("expected method GET, got %s", req.Method)
	}
	if req.URL != "https://httpbin.org/get" {
		t.Errorf("expected URL https://httpbin.org/get, got %s", req.URL)
	}
}

func TestParseRequestJSON(t *testing.T) {
	input := `{"method": "POST", "url": "https://example.com/api"}`
	req, err := ParseRequest(input)
	if err != nil {
		t.Fatalf("ParseRequest failed: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("expected method POST, got %s", req.Method)
	}
}

func TestLoadCollection(t *testing.T) {
	content := `
name: Test Collection
description: A test collection
base_url: https://httpbin.org
requests:
  - name: Get
    method: GET
    url: /get
  - name: Post
    method: POST
    url: /post
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "collection.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	col, err := LoadCollection(path)
	if err != nil {
		t.Fatalf("LoadCollection failed: %v", err)
	}

	if col.Name != "Test Collection" {
		t.Errorf("expected name 'Test Collection', got %s", col.Name)
	}
	if len(col.Requests) != 2 {
		t.Errorf("expected 2 requests, got %d", len(col.Requests))
	}
}

func TestLoadEnvironment(t *testing.T) {
	content := `
name: dev
variables:
  API_KEY: test123
  BASE_URL: http://localhost:8080
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "env.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	env, err := LoadEnvironment(path)
	if err != nil {
		t.Fatalf("LoadEnvironment failed: %v", err)
	}

	if env.Name != "dev" {
		t.Errorf("expected name 'dev', got %s", env.Name)
	}
	if env.Variables["API_KEY"] != "test123" {
		t.Errorf("expected API_KEY=test123, got %s", env.Variables["API_KEY"])
	}
}
