package turso

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

const (
	databaseEndpoint = "v1/organizations/%s/databases"
	maxNameLength    = 32
	regexName        = "^[a-z0-9-]+$"
)

// DatabaseService is the interface for the Turso API database endpoint
type DatabaseService service

type databaseService interface {
	// ListDatabases lists all databases in the organization
	ListDatabases(ctx context.Context) (*ListDatabaseResponse, error)
	// CreateDatabase creates a new database
	CreateDatabase(ctx context.Context, req CreateDatabaseRequest) (*CreateDatabaseResponse, error)
	// GetDatabase gets a database by name
	GetDatabase(ctx context.Context, dbName string) (*GetDatabaseResponse, error)
	// DeleteDatabase deletes a database by name
	DeleteDatabase(ctx context.Context, dbName string) (*DeleteDatabaseResponse, error)
}

// Database is the struct for the Turso Database object
type Database struct {
	// Name is the name of the database
	Name string `json:"Name"`
	// DatabaseID is the ID of the database
	DatabaseID string `json:"DbId"`
	// Hostname is the hostname of the database`
	Hostname string `json:"Hostname"` // this is in the response twice, once with a capital H and once with a lowercase h
	// IsSchema is this database controls other child databases then this will be true
	IsSchema bool `json:"is_schema"`
	// Schema is the name of the parent database that owns the schema for this database
	Schema string `json:"schema"`
	// BlockedReads is true if reads are blocked
	BlockReads bool `json:"block_reads"`
	// BlockedWrites is true if writes are blocked
	BlockWrites bool `json:"block_writes"`
	// AllowAttach is true if the database allows attachments of a child database
	AllowAttach bool `json:"allow_attach"`
	// Regions is a list of regions the database is available in
	Regions []string `json:"regions"`
	// PrimaryRegion is the primary region for the database
	PrimaryRegion string `json:"primaryRegion"`
	// Type is the type of the database
	Type string `json:"type"`
	// Version is the version of libsql used by the database
	Version string `json:"version"`
	// Group is the group the database is in
	Group string `json:"group"`
	// Sleeping is true if the database is sleeping
	Sleeping bool `json:"sleeping"`
}

// CreateDatabase is the struct for the Turso API database create request
type CreateDatabase struct {
	// DatabaseID is the ID of the database
	DatabaseID string `json:"DbId"`
	// Name is the name of the database
	Name string `json:"Name"`
	// Hostname is the hostname of the database
	Hostname string `json:"Hostname"`
	// IssuedCertCount is the number of certificates issued
	IssuedCertCount int `json:"IssuedCertCount"`
	// IssuedCertLimit is the limit of certificates that can be issued
	IssuedCertLimit int `json:"IssuedCertLimit"`
}

// ListDatabaseResponse is the struct for the Turso API database list response
type ListDatabaseResponse struct {
	Databases []*Database `json:"databases"`
}

// GetDatabaseResponse is the struct for the Turso API database get response
type GetDatabaseResponse struct {
	Database *Database `json:"database"`
}

// CreateDatabaseResponse is the struct for the Turso API database create response
type CreateDatabaseResponse struct {
	Database CreateDatabase `json:"database"`
}

// DeleteDatabaseResponse is the struct for the Turso API database delete response
type DeleteDatabaseResponse struct {
	Database string `json:"database"`
}

// CreateDatabaseRequest is the struct for the Turso API database create request
type CreateDatabaseRequest struct {
	// Group is the group the database is in
	Group string `json:"group"`
	// IsSchema is this database controls other child databases then this will be true
	IsSchema bool `json:"is_schema"`
	// Name is the name of the database
	// Must contain only lowercase letters, numbers, dashes. No longer than 32 characters.
	Name string `json:"name"`
}

// getDatabaseEndpoint returns the endpoint for the Turso API database service
func getDatabaseEndpoint(baseURL, orgName string) string {
	dbEndpoint := fmt.Sprintf(databaseEndpoint, orgName)
	return fmt.Sprintf("%s/%s", baseURL, dbEndpoint)
}

// CreateDatabase satisfies the databaseService interface
func (s *DatabaseService) CreateDatabase(ctx context.Context, db CreateDatabaseRequest) (*CreateDatabaseResponse, error) {
	// Sanitize the database name
	if err := validateDatabaseName(db.Name); err != nil {
		return nil, err
	}

	// Create the database
	endpoint := getDatabaseEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName)

	resp, err := s.client.DoRequest(ctx, http.MethodPost, endpoint, db)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode the response
	var out CreateDatabaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("database", "creating", resp.StatusCode)
	}

	return &out, nil
}

// ListDatabases satisfies the databaseService interface
func (s *DatabaseService) ListDatabases(ctx context.Context) (*ListDatabaseResponse, error) {
	endpoint := getDatabaseEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName)

	resp, err := s.client.DoRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var out ListDatabaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("databases", "listing", resp.StatusCode)
	}

	return &out, nil
}

// GetDatabase satisfies the databaseService interface
func (s *DatabaseService) GetDatabase(ctx context.Context, dbName string) (*GetDatabaseResponse, error) {
	// get endpoint and append the database name
	endpoint := getDatabaseEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName)
	endpoint = fmt.Sprintf("%s/%s", endpoint, dbName)

	resp, err := s.client.DoRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var out *GetDatabaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("database", "getting", resp.StatusCode)
	}

	return out, nil
}

// DeleteDatabase satisfies the databaseService interface
func (s *DatabaseService) DeleteDatabase(ctx context.Context, dbName string) (*DeleteDatabaseResponse, error) {
	// Delete the database
	endpoint := getDatabaseEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName)
	endpoint = fmt.Sprintf("%s/%s", endpoint, dbName)

	resp, err := s.client.DoRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode the response
	var out DeleteDatabaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("database", "deleting", resp.StatusCode)
	}

	return &out, nil
}

// validateDatabaseName validates the database name to ensure it meets the requirements set by the Turso API
func validateDatabaseName(name string) error {
	match, err := regexp.MatchString(regexName, name)
	if err != nil {
		return err
	}

	if !match {
		return ErrInvalidDatabaseName
	}

	if len(name) > maxNameLength {
		return ErrInvalidDatabaseName
	}

	return nil
}
