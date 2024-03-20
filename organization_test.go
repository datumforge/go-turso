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

func TestListOrganizations(t *testing.T) {
	body := `[{"name":"personal","slug":"moewlemonade","type":"personal","plan_id":"starter","overages":false,"blocked_reads":false,"blocked_writes":false,"plan_timeline":"","memory":0},{"name":"meow","slug":"meow","type":"team","plan_id":"scaler","overages":false,"blocked_reads":false,"blocked_writes":false,"plan_timeline":"monthly","memory":0}]`
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

	orgService := OrganizationService{client: client}
	resp, err := orgService.ListOrganizations(context.Background())
	require.NoError(t, err)
	assert.Len(t, *resp, 2)
}
