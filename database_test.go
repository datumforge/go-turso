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

func TestListDatabases(t *testing.T) {
	databaseService := newMockDatabaseService()

	resp, err := databaseService.ListDatabases(context.Background())
	require.NoError(t, err)
	assert.Len(t, resp.Databases, 1)
}
func TestGetDatabase(t *testing.T) {
	databaseService := newMockDatabaseService()

	resp, err := databaseService.GetDatabase(context.Background(), "my-db")
	require.NoError(t, err)
	assert.Equal(t, resp.Database.Name, "my-db")
}

func TestDeleteDatabase(t *testing.T) {
	databaseService := newMockDatabaseService()

	resp, err := databaseService.DeleteDatabase(context.Background(), "my-db")
	require.NoError(t, err)
	assert.Equal(t, resp.Database, "my-db")
}

func TestCreateDatabase(t *testing.T) {
	body := `{"database":{"DbId":"0eb771dd-6906-11ee-8553-eaa7715aeaf2","Hostname":"[databaseName]-[organizationName].turso.io","Name":"my-db"}}`
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
	databaseService := DatabaseService{client: client}
	req := CreateDatabaseRequest{
		Name: "my-db",
	}

	resp, err := databaseService.CreateDatabase(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, resp.Database.Name, "my-db")

	// test error
	req = CreateDatabaseRequest{
		Name: "myAWESOMEdb",
	}

	resp, err = databaseService.CreateDatabase(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestValidateDatabaseName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "happy path, simple name",
			input:     "mydatabase-123",
			expectErr: false,
		},
		{
			name:      "name with uppercase",
			input:     "ORG-123ABC",
			expectErr: true,
		},
		{
			name:      "long name",
			input:     "MECPJTpHtBEyUNBAujXw6mxCjN4ARLPJ3",
			expectErr: true,
		},
	}

	for _, test := range tests {
		result := validateDatabaseName(test.input)
		if test.expectErr {
			require.Error(t, result)
		} else {
			require.NoError(t, result)
		}
	}
}
