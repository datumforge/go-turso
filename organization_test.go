package turso

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListOrganizations(t *testing.T) {
	orgService := newMockOrganizationService()

	resp, err := orgService.ListOrganizations(context.Background())
	require.NoError(t, err)
	assert.Len(t, *resp, 1)
}
