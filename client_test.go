package reader

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "valid-token",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Errorf("NewClient() returned nil client")
			}
		})
	}
}

func TestClientError(t *testing.T) {
	err := &ClientError{
		Type:    "test_error",
		Message: "test message",
	}
	expected := "test_error: test message"
	if err.Error() != expected {
		t.Errorf("expected error %s, got %s", expected, err.Error())
	}
}

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: 400,
		Message:    "Bad Request",
		Details:    map[string]interface{}{"field": "error"},
	}
	expected := "API error (status 400): Bad Request"
	if err.Error() != expected {
		t.Errorf("expected error %s, got %s", expected, err.Error())
	}
}