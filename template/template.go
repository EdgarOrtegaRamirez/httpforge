// Package template handles variable substitution in request configurations.
package template

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Engine handles template variable substitution.
type Engine struct {
	vars map[string]string
}

// NewEngine creates a new template engine with the given variables.
func NewEngine(vars map[string]string) *Engine {
	e := &Engine{
		vars: make(map[string]string),
	}
	if vars != nil {
		e.Merge(vars)
	}
	return e
}

// Merge adds variables to the engine, with higher priority taking precedence.
func (e *Engine) Merge(vars map[string]string) {
	for k, v := range vars {
		e.vars[k] = v
	}
}

// LoadEnv loads environment variables with the given prefix.
func (e *Engine) LoadEnv(prefix string) {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			if prefix != "" {
				if strings.HasPrefix(key, prefix) {
					key = strings.TrimPrefix(key, prefix)
				} else {
					continue
				}
			}
			e.vars[key] = parts[1]
		}
	}
}

// Render replaces ${var} placeholders in the input string.
func (e *Engine) Render(input string) string {
	return varPattern.ReplaceAllStringFunc(input, func(match string) string {
		varName := varPattern.FindStringSubmatch(match)[1]

		// Check direct variables first
		if val, ok := e.vars[varName]; ok {
			return val
		}

		// Check environment variables
		if val, ok := os.LookupEnv(varName); ok {
			return val
		}

		// Return original placeholder if not found
		return match
	})
}

// RenderMap renders all values in a map.
func (e *Engine) RenderMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[e.Render(k)] = e.Render(v)
	}
	return result
}

// ExtractVarName extracts the variable name from a ${var} pattern.
func ExtractVarName(pattern string) string {
	matches := varPattern.FindStringSubmatch(pattern)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// HasVariables checks if a string contains template variables.
func HasVariables(s string) bool {
	return varPattern.MatchString(s)
}

// ListVariables returns all variable names referenced in a string.
func ListVariables(s string) []string {
	matches := varPattern.FindAllStringSubmatch(s, -1)
	var vars []string
	seen := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 && !seen[m[1]] {
			vars = append(vars, m[1])
			seen[m[1]] = true
		}
	}
	return vars
}

// Validate checks that all referenced variables are available.
func (e *Engine) Validate(input string) []string {
	var missing []string
	for _, varName := range ListVariables(input) {
		if _, ok := e.vars[varName]; !ok {
			if _, ok := os.LookupEnv(varName); !ok {
				missing = append(missing, varName)
			}
		}
	}
	return missing
}

// String returns a string representation of the engine's variables (for debugging).
func (e *Engine) String() string {
	var parts []string
	for k, v := range e.vars {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ", ")
}

// Vars returns the current variable map.
func (e *Engine) Vars() map[string]string {
	return e.vars
}
