// Package models defines the core data structures for HttpForge.
package models

import (
	"time"
)

// HTTPMethod represents an HTTP request method.
type HTTPMethod string

const (
	MethodGET     HTTPMethod = "GET"
	MethodPOST    HTTPMethod = "POST"
	MethodPUT     HTTPMethod = "PUT"
	MethodPATCH   HTTPMethod = "PATCH"
	MethodDELETE  HTTPMethod = "DELETE"
	MethodHEAD    HTTPMethod = "HEAD"
	MethodOPTIONS HTTPMethod = "OPTIONS"
)

// Request represents an HTTP request configuration.
type Request struct {
	Name        string            `yaml:"name" json:"name"`
	Method      HTTPMethod        `yaml:"method" json:"method"`
	URL         string            `yaml:"url" json:"url"`
	Headers     map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Body        string            `yaml:"body,omitempty" json:"body,omitempty"`
	BodyType    string            `yaml:"body_type,omitempty" json:"body_type,omitempty"` // json, form, raw
	QueryParams map[string]string `yaml:"query_params,omitempty" json:"query_params,omitempty"`
	Timeout     time.Duration     `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Assertions  []Assertion       `yaml:"assertions,omitempty" json:"assertions,omitempty"`
	Variables   map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"` // Extract response data
}

// Response represents an HTTP response.
type Response struct {
	StatusCode    int               `json:"status_code"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"body"`
	Timing        time.Duration     `json:"timing"`
	Size          int               `json:"size"`
	Request       *Request          `json:"request,omitempty"`
	ExtractedVars map[string]string `json:"extracted_vars,omitempty"`
}

// Environment represents a set of variables for request templating.
type Environment struct {
	Name      string            `yaml:"name" json:"name"`
	Variables map[string]string `yaml:"variables" json:"variables"`
}

// Collection represents a group of requests.
type Collection struct {
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	BaseURL     string            `yaml:"base_url,omitempty" json:"base_url,omitempty"`
	Variables   map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"`
	Requests    []Request         `yaml:"requests" json:"requests"`
}

// Assertion represents a response assertion.
type Assertion struct {
	Field    string `yaml:"field" json:"field"`                   // status_code, header, body, json_path
	Path     string `yaml:"path,omitempty" json:"path,omitempty"` // JSON path for json_path field
	Op       string `yaml:"op" json:"op"`                         // eq, ne, gt, lt, contains, matches, exists
	Value    string `yaml:"value" json:"value"`
	Negative bool   `yaml:"negative,omitempty" json:"negative,omitempty"` // assert NOT
}

// AssertionResult represents the result of an assertion check.
type AssertionResult struct {
	Assertion Assertion `json:"assertion"`
	Passed    bool      `json:"passed"`
	Message   string    `json:"message"`
}

// RequestConfig is the unified request configuration that can be loaded from file.
type RequestConfig struct {
	Request     Request      `yaml:"request" json:"request"`
	Environment *Environment `yaml:"environment,omitempty" json:"environment,omitempty"`
}

// ExportFormat represents the format for code export.
type ExportFormat string

const (
	ExportCurl   ExportFormat = "curl"
	ExportPython ExportFormat = "python"
	ExportJS     ExportFormat = "javascript"
)

// OutputFormat represents the output display format.
type OutputFormat string

const (
	OutputText OutputFormat = "text"
	OutputJSON OutputFormat = "json"
	OutputCSV  OutputFormat = "csv"
)
