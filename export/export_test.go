package export

import (
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/httpforge/models"
)

func TestToCurlGET(t *testing.T) {
	req := &models.Request{
		Method: models.MethodGET,
		URL:    "https://httpbin.org/get",
	}

	result := ToCurl(req)
	if !strings.Contains(result, "curl") {
		t.Error("expected curl command to contain 'curl'")
	}
	if !strings.Contains(result, "https://httpbin.org/get") {
		t.Error("expected curl command to contain URL")
	}
}

func TestToCurlPOST(t *testing.T) {
	req := &models.Request{
		Method: models.MethodPOST,
		URL:    "https://httpbin.org/post",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:     `{"key": "value"}`,
		BodyType: "json",
	}

	result := ToCurl(req)
	if !strings.Contains(result, "POST") {
		t.Error("expected curl command to contain POST method")
	}
	if !strings.Contains(result, "Content-Type") {
		t.Error("expected curl command to contain Content-Type header")
	}
	if !strings.Contains(result, `{"key": "value"}`) {
		t.Error("expected curl command to contain body")
	}
}

func TestToCurlWithQueryParams(t *testing.T) {
	req := &models.Request{
		Method: models.MethodGET,
		URL:    "https://httpbin.org/get",
		QueryParams: map[string]string{
			"foo": "bar",
			"baz": "qux",
		},
	}

	result := ToCurl(req)
	if !strings.Contains(result, "foo=bar") {
		t.Error("expected curl command to contain query params")
	}
}

func TestToPython(t *testing.T) {
	req := &models.Request{
		Method: models.MethodGET,
		URL:    "https://httpbin.org/get",
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	result := ToPython(req)
	if !strings.Contains(result, "import requests") {
		t.Error("expected Python code to contain 'import requests'")
	}
	if !strings.Contains(result, "requests.get") {
		t.Error("expected Python code to contain 'requests.get'")
	}
}

func TestToPythonPOST(t *testing.T) {
	req := &models.Request{
		Method:   models.MethodPOST,
		URL:      "https://httpbin.org/post",
		BodyType: "json",
		Body:     `{"key": "value"}`,
	}

	result := ToPython(req)
	if !strings.Contains(result, "requests.post") {
		t.Error("expected Python code to contain 'requests.post'")
	}
	if !strings.Contains(result, "json=") {
		t.Error("expected Python code to contain 'json='")
	}
}

func TestToJavaScript(t *testing.T) {
	req := &models.Request{
		Method: models.MethodGET,
		URL:    "https://httpbin.org/get",
	}

	result := ToJavaScript(req)
	if !strings.Contains(result, "fetch") {
		t.Error("expected JavaScript code to contain 'fetch'")
	}
	if !strings.Contains(result, "async function") {
		t.Error("expected JavaScript code to contain 'async function'")
	}
}

func TestExportCurl(t *testing.T) {
	req := &models.Request{
		Method: models.MethodGET,
		URL:    "https://httpbin.org/get",
	}

	result := Export(req, models.ExportCurl)
	if !strings.Contains(result, "curl") {
		t.Error("expected curl export")
	}
}

func TestExportPython(t *testing.T) {
	req := &models.Request{
		Method: models.MethodGET,
		URL:    "https://httpbin.org/get",
	}

	result := Export(req, models.ExportPython)
	if !strings.Contains(result, "requests") {
		t.Error("expected Python export")
	}
}

func TestExportJavaScript(t *testing.T) {
	req := &models.Request{
		Method: models.MethodGET,
		URL:    "https://httpbin.org/get",
	}

	result := Export(req, models.ExportJS)
	if !strings.Contains(result, "fetch") {
		t.Error("expected JavaScript export")
	}
}
