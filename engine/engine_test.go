package engine

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/httpforge/models"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("NewClient returned nil")
	}
	if client.httpClient == nil {
		t.Error("httpClient is nil")
	}
}

func TestExecuteWithInvalidURL(t *testing.T) {
	client := NewClient()

	req := &models.Request{
		Method: models.MethodGET,
		URL:    "http://localhost:1", // Invalid URL that will fail
	}

	_, err := client.Execute(req)
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestExecuteWithInvalidMethod(t *testing.T) {
	client := NewClient()

	req := &models.Request{
		Method: "INVALID",
		URL:    "http://localhost:1",
	}

	_, err := client.Execute(req)
	if err == nil {
		t.Error("expected error for invalid method")
	}
}
