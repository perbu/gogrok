package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockResponseWriter is a simple implementation of http.ResponseWriter for testing
type mockResponseWriter struct {
	headers    http.Header
	statusCode int
	body       []byte
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		headers: make(http.Header),
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.headers
}

func (m *mockResponseWriter) Write(body []byte) (int, error) {
	m.body = append(m.body, body...)
	return len(body), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func TestResponseObserverWrite(t *testing.T) {
	tests := []struct {
		name           string
		writeHeader    bool
		expectedStatus int
		writeData      string
		expectedBytes  int64
	}{
		{
			name:           "write without explicit header",
			writeHeader:    false,
			expectedStatus: http.StatusOK,
			writeData:      "test data",
			expectedBytes:  9, // len("test data")
		},
		{
			name:           "write with explicit header",
			writeHeader:    true,
			expectedStatus: http.StatusCreated,
			writeData:      "test data",
			expectedBytes:  9, // len("test data")
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRW := newMockResponseWriter()
			observer := &responseObserver{
				ResponseWriter: mockRW,
			}

			if tt.writeHeader {
				observer.WriteHeader(http.StatusCreated)
			}

			n, err := observer.Write([]byte(tt.writeData))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if n != len(tt.writeData) {
				t.Errorf("expected to write %d bytes, got %d", len(tt.writeData), n)
			}

			if observer.status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, observer.status)
			}
			if observer.written != tt.expectedBytes {
				t.Errorf("expected written bytes %d, got %d", tt.expectedBytes, observer.written)
			}
			if !observer.wroteHeader {
				t.Errorf("expected wroteHeader to be true")
			}
		})
	}
}

func TestResponseObserverWriteHeader(t *testing.T) {
	tests := []struct {
		name           string
		callTwice      bool
		firstStatus    int
		secondStatus   int
		expectedStatus int
	}{
		{
			name:           "write header once",
			callTwice:      false,
			firstStatus:    http.StatusCreated,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "write header twice",
			callTwice:      true,
			firstStatus:    http.StatusCreated,
			secondStatus:   http.StatusBadRequest,
			expectedStatus: http.StatusCreated, // First status should be preserved
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRW := newMockResponseWriter()
			observer := &responseObserver{
				ResponseWriter: mockRW,
			}

			observer.WriteHeader(tt.firstStatus)
			
			if tt.callTwice {
				observer.WriteHeader(tt.secondStatus)
			}

			if observer.status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, observer.status)
			}
			if !observer.wroteHeader {
				t.Errorf("expected wroteHeader to be true")
			}
			if mockRW.statusCode != tt.expectedStatus {
				t.Errorf("expected mockRW status %d, got %d", tt.expectedStatus, mockRW.statusCode)
			}
		})
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Create a simple handler that returns a 200 status code
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})

	// Wrap the handler with the logging middleware
	wrappedHandler := loggingMiddleware(handler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Call the wrapped handler
	wrappedHandler.ServeHTTP(rec, req)

	// Check the response
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "test" {
		t.Errorf("expected body %q, got %q", "test", rec.Body.String())
	}
}