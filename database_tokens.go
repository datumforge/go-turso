package turso

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xhit/go-str2duration/v2"
)

const (
	databaseTokensEndpoint = "v1/organizations/%s/databases/%s/auth/tokens"
	FullAccess             = "full-access"
	ReadOnly               = "read-only"
	DefaultExpiration      = "never"
)

var validAuthorization = []string{FullAccess, ReadOnly}

// DatabaseTokensService is the interface for the Turso API database tokens service
type DatabaseTokensService service

type databaseTokensService interface {
	// CreateDatabaseToken creates a new database token
	CreateDatabaseToken(ctx context.Context, req CreateDatabaseTokenRequest) (*CreateDatabaseTokenResponse, error)
}

// CreateDatabaseTokenRequest is the struct for the Turso API database token create request
type CreateDatabaseTokenRequest struct {
	// DatabaseName is the name of the database
	DatabaseName string
	// Expiration is the expiration time for the token
	Expiration string
	// Permissions is the permissions for the token
	Authorization string
	// ReadAttach permission for the token
	AttachPermissions Permissions `json:"permissions"`
}

type Permissions struct {
	ReadAttach struct {
		Database []string `json:"database"`
	} `json:"read_attach"`
}

// CreateDatabaseTokenResponse is the struct for the Turso API database token create response
type CreateDatabaseTokenResponse struct {
	JWT string `json:"jwt"`
}

// InvalidateDatabaseTokenRequest is the struct for the Turso API database token invalidate request
type InvalidateDatabaseTokenRequest struct {
	// DatabaseName is the name of the database
	DatabaseName string `json:"database_name"`
}

// getDatabaseTokensEndpoint returns the endpoint for the Turso API database token service
func getDatabaseTokensEndpoint(baseURL, orgName, dbName string) string {
	dbEndpoint := fmt.Sprintf(databaseTokensEndpoint, orgName, dbName)
	return fmt.Sprintf("%s/%s", baseURL, dbEndpoint)
}

// CreateDatabaseToken satisfies the databaseTokensService interface
func (s *DatabaseTokensService) CreateDatabaseToken(ctx context.Context, req CreateDatabaseTokenRequest) (*CreateDatabaseTokenResponse, error) {
	if err := validateDatabaseTokenRequest(req); err != nil {
		return nil, err
	}

	endpoint := getDatabaseTokensEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName, req.DatabaseName)
	endpoint = fmt.Sprintf("%s?expiration=%s&authorization=%s", endpoint, req.Expiration, req.Authorization)

	resp, err := s.client.DoRequest(ctx, http.MethodPost, endpoint, req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var out CreateDatabaseTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("database token", "creating", resp.StatusCode)
	}

	return &out, nil
}

// validateDatabaseTokenRequest ensures the authorization and expiration are valid
// in the given request before making the API call
func validateDatabaseTokenRequest(req CreateDatabaseTokenRequest) error {
	if !isValidExpiration(req.Expiration) {
		return ErrExpirationInvalid
	}

	if !isValidAuthorization(req.Authorization) {
		return ErrAuthorizationInvalid
	}

	return nil
}

// IsValidExpiration checks if the expiration is valid
func isValidExpiration(expiration string) bool {
	// check for empty fields first
	if expiration == "" {
		return false
	}

	if expiration == DefaultExpiration {
		return true
	}

	if _, err := str2duration.ParseDuration(expiration); err != nil {
		return false
	}

	return true
}

// IsValidAuthorization checks if the authorization is valid
func isValidAuthorization(authorization string) bool {
	// check for empty fields first
	if authorization == "" {
		return false
	}

	// check for valid authorization
	for _, v := range validAuthorization {
		if v == authorization {
			return true
		}
	}

	return false
}
