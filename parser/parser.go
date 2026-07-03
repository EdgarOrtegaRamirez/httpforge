// Package parser handles loading request configurations from YAML/JSON files.
package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/httpforge/models"
	"gopkg.in/yaml.v3"
)

// LoadRequestFile loads a request configuration from a YAML or JSON file.
func LoadRequestFile(path string) (*models.RequestConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	var config models.RequestConfig

	if strings.HasSuffix(path, ".json") {
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("parsing JSON: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("parsing YAML: %w", err)
		}
	}

	// Set defaults
	if config.Request.Method == "" {
		config.Request.Method = models.MethodGET
	}
	if config.Request.Timeout == 0 {
		config.Request.Timeout = 30_000_000_000 // 30 seconds as time.Duration (nanoseconds)
	}

	return &config, nil
}

// LoadCollection loads a collection of requests from a YAML or JSON file.
func LoadCollection(path string) (*models.Collection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	var col models.Collection

	if strings.HasSuffix(path, ".json") {
		if err := json.Unmarshal(data, &col); err != nil {
			return nil, fmt.Errorf("parsing JSON: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &col); err != nil {
			return nil, fmt.Errorf("parsing YAML: %w", err)
		}
	}

	return &col, nil
}

// LoadEnvironment loads environment variables from a YAML/JSON file.
func LoadEnvironment(path string) (*models.Environment, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	var env models.Environment

	if strings.HasSuffix(path, ".json") {
		if err := json.Unmarshal(data, &env); err != nil {
			return nil, fmt.Errorf("parsing JSON: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &env); err != nil {
			return nil, fmt.Errorf("parsing YAML: %w", err)
		}
	}

	return &env, nil
}

// ParseRequest parses a request from a string (URL shorthand or full config).
func ParseRequest(input string) (*models.Request, error) {
	// If it looks like a URL, create a simple GET request
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return &models.Request{
			Method: models.MethodGET,
			URL:    input,
		}, nil
	}

	// Try to parse as inline JSON
	var req models.Request
	if err := json.Unmarshal([]byte(input), &req); err == nil {
		if req.Method == "" {
			req.Method = models.MethodGET
		}
		return &req, nil
	}

	// Try to parse as inline YAML
	if err := yaml.Unmarshal([]byte(input), &req); err == nil {
		if req.Method == "" {
			req.Method = models.MethodGET
		}
		return &req, nil
	}

	return nil, fmt.Errorf("cannot parse input: %s", input)
}
