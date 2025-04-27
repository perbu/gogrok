package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFoundHandler(t *testing.T) {
	// Create a test request
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rec := httptest.NewRecorder()

	// Call the handler
	notFoundHandler(rec, req)

	// Check the response
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, rec.Code)
	}

	// Check the response body
	expectedBody := "not found\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestMakeStaticHandler(t *testing.T) {
	// Create a test request
	req := httptest.NewRequest("GET", "/styles.css", nil)
	rec := httptest.NewRecorder()

	// Create a handler for a file that exists in the embedded filesystem
	handler := makeStaticHandler("assets/styles.css")

	// Call the handler
	handler(rec, req)

	// Check the response
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// We can't easily check the exact response body since it depends on the embedded file,
	// but we can at least verify that some response was written
	if rec.Body.Len() == 0 {
		t.Errorf("expected non-empty response body")
	}

	// Test with a non-existent file
	req = httptest.NewRequest("GET", "/nonexistent.css", nil)
	rec = httptest.NewRecorder()

	// Create a handler for a file that doesn't exist
	handler = makeStaticHandler("assets/nonexistent.css")

	// Call the handler
	handler(rec, req)

	// Check the response
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, rec.Code)
	}

	// Check the response body
	expectedBody := "file not found\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}
