// HttpForge is a lightweight, config-driven HTTP API testing toolkit.
//
// It allows you to build, execute, and validate HTTP requests from
// YAML/JSON configuration files. Features include:
//   - YAML/JSON request configs
//   - Template variable substitution
//   - Response assertions
//   - Code export (curl, Python, JavaScript)
//   - Multiple output formats (text, JSON, CSV)
//   - Environment variable management
package main

import (
	"github.com/EdgarOrtegaRamirez/httpforge/cmd"
)

func main() {
	cmd.Execute()
}
