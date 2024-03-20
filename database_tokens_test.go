package turso

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDatabaseToken(t *testing.T) {
	body := `{"jwt": "areallylongstringjwtgoeshere"}`
	client := &Client{
		cfg: &Config{
			BaseURL: "http://localhost",
		},
		client: &MockHTTPRequestDoer{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(body))),
			},
		},
	}

	// happy path
	databaseTokenService := DatabaseTokensService{client: client}
	req := CreateDatabaseTokenRequest{
		DatabaseName:  "my-db",
		Expiration:    "1h30m",
		Authorization: FullAccess,
	}

	resp, err := databaseTokenService.CreateDatabaseToken(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp.JWT)

	// test errors
	req = CreateDatabaseTokenRequest{
		DatabaseName:  "my-db",
		Authorization: FullAccess,
	}

	resp, err = databaseTokenService.CreateDatabaseToken(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	req = CreateDatabaseTokenRequest{
		DatabaseName:  "my-db",
		Expiration:    "1h30m",
		Authorization: "invalid",
	}

	resp, err = databaseTokenService.CreateDatabaseToken(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestValidateDatabaseTokenRequest(t *testing.T) {
	tests := []struct {
		name    string
		request CreateDatabaseTokenRequest
		wantErr error
	}{
		{
			name: "Valid request, full access",
			request: CreateDatabaseTokenRequest{
				Expiration:    "never",
				Authorization: "full-access",
			},
			wantErr: nil,
		},
		{
			name: "Valid request, read only",
			request: CreateDatabaseTokenRequest{
				Expiration:    "never",
				Authorization: "read-only",
			},
			wantErr: nil,
		},
		{
			name: "Valid request, with duration",
			request: CreateDatabaseTokenRequest{
				Expiration:    "12w",
				Authorization: "read-only",
			},
			wantErr: nil,
		},
		{
			name: "Missing expiration",
			request: CreateDatabaseTokenRequest{
				Expiration:    "",
				Authorization: "read-only",
			},
			wantErr: ErrExpirationInvalid,
		},
		{
			name: "Invalid authorization",
			request: CreateDatabaseTokenRequest{
				Expiration:    "never",
				Authorization: "invalid",
			},
			wantErr: ErrAuthorizationInvalid,
		},
		{
			name: "Invalid expiration",
			request: CreateDatabaseTokenRequest{
				Expiration:    "2030-01-01",
				Authorization: "invalid",
			},
			wantErr: ErrExpirationInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDatabaseTokenRequest(tt.request)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}

			require.NoError(t, err)
		})
	}
}
