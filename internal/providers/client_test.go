package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetLatestVersion(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		statusCode  int
		expectError bool
		wantVersion string
	}{
		{
			name:        "Success",
			response:    `{"version": "5.30.0"}`,
			statusCode:  http.StatusOK,
			expectError: false,
			wantVersion: "5.30.0",
		},
		{
			name:        "NotFound",
			response:    "",
			statusCode:  http.StatusNotFound,
			expectError: true,
			wantVersion: "",
		},
		{
			name:        "BadJSON",
			response:    `{invalid json}`,
			statusCode:  http.StatusOK,
			expectError: true,
			wantVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// implementation detail: Client needs to point to test server
			c := &Client{
				BaseURL:    server.URL,
				HTTPClient: server.Client(),
			}

			// We pass "source" but with our modified BaseURL, it effectively ignores the real registry
			got, err := c.GetLatestVersion("hashicorp/aws")

			if (err != nil) != tt.expectError {
				t.Errorf("GetLatestVersion() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if got != tt.wantVersion {
				t.Errorf("GetLatestVersion() = %v, want %v", got, tt.wantVersion)
			}
		})
	}
}
