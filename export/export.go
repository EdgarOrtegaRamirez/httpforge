// Package export handles generating code representations of HTTP requests.
package export

import (
	"fmt"
	"strings"

	"github.com/EdgarOrtegaRamirez/httpforge/models"
)

// ToCurl generates a curl command from a request.
func ToCurl(req *models.Request) string {
	var parts []string
	parts = append(parts, "curl")

	// Method
	if req.Method != models.MethodGET {
		parts = append(parts, "-X", string(req.Method))
	}

	// Headers
	for key, value := range req.Headers {
		parts = append(parts, "-H", fmt.Sprintf("'%s: %s'", key, value))
	}

	// Body
	if req.Body != "" {
		// Escape single quotes in body
		escaped := strings.ReplaceAll(req.Body, "'", "'\\''")
		parts = append(parts, "-d", fmt.Sprintf("'%s'", escaped))
	}

	// URL with query params
	url := req.URL
	if len(req.QueryParams) > 0 {
		var params []string
		for key, value := range req.QueryParams {
			params = append(params, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(params, "&")
	}
	parts = append(parts, fmt.Sprintf("'%s'", url))

	return strings.Join(parts, " \\\n  ")
}

// ToPython generates Python requests code from a request.
func ToPython(req *models.Request) string {
	var lines []string
	lines = append(lines, "import requests")

	// Headers
	if len(req.Headers) > 0 {
		lines = append(lines, "")
		lines = append(lines, "headers = {")
		for key, value := range req.Headers {
			lines = append(lines, fmt.Sprintf("    \"%s\": \"%s\",", key, value))
		}
		lines = append(lines, "}")
	}

	// URL with query params
	url := req.URL
	if len(req.QueryParams) > 0 {
		var params []string
		for key, value := range req.QueryParams {
			params = append(params, fmt.Sprintf("\"%s\": \"%s\"", key, value))
		}
		url += "?" + strings.Join(params, "&")
	}

	// Make request
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("response = requests.%s(", strings.ToLower(string(req.Method))))

	if len(req.Headers) > 0 {
		lines = append(lines, fmt.Sprintf("    \"%s\",", url))
		lines = append(lines, "    headers=headers,")
	} else {
		lines = append(lines, fmt.Sprintf("    \"%s\",", url))
	}

	if req.Body != "" {
		if req.BodyType == "json" {
			lines = append(lines, "    json="+req.Body+",")
		} else {
			lines = append(lines, "    data=\""+req.Body+"\",")
		}
	}

	lines = append(lines, ")")
	lines = append(lines, "")
	lines = append(lines, "print(f\"Status: {response.status_code}\")")
	lines = append(lines, "print(f\"Body: {response.text}\")")

	return strings.Join(lines, "\n")
}

// ToJavaScript generates JavaScript fetch code from a request.
func ToJavaScript(req *models.Request) string {
	var lines []string
	lines = append(lines, "async function makeRequest() {")

	// URL
	url := req.URL
	if len(req.QueryParams) > 0 {
		var params []string
		for key, value := range req.QueryParams {
			params = append(params, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(params, "&")
	}

	lines = append(lines, fmt.Sprintf("  const response = await fetch(\"%s\", {", url))

	// Method
	if req.Method != models.MethodGET {
		lines = append(lines, fmt.Sprintf("    method: \"%s\",", req.Method))
	}

	// Headers
	if len(req.Headers) > 0 {
		lines = append(lines, "    headers: {")
		for key, value := range req.Headers {
			lines = append(lines, fmt.Sprintf("      \"%s\": \"%s\",", key, value))
		}
		lines = append(lines, "    },")
	}

	// Body
	if req.Body != "" {
		if req.BodyType == "json" {
			lines = append(lines, "    body: JSON.stringify("+req.Body+"),")
		} else {
			lines = append(lines, fmt.Sprintf("    body: \"%s\",", req.Body))
		}
	}

	lines = append(lines, "  });")
	lines = append(lines, "")
	lines = append(lines, "  const data = await response.json();")
	lines = append(lines, "  console.log(\"Status:\", response.status);")
	lines = append(lines, "  console.log(\"Body:\", data);")
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

// Export generates code in the specified format.
func Export(req *models.Request, format models.ExportFormat) string {
	switch format {
	case models.ExportCurl:
		return ToCurl(req)
	case models.ExportPython:
		return ToPython(req)
	case models.ExportJS:
		return ToJavaScript(req)
	default:
		return ToCurl(req)
	}
}
