package turso

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	organizationEndpoint = "v1/organizations"
)

type OrganizationService service

type organizationService interface {
	// ListOrganizations lists all organizations for the authorized user
	ListOrganizations(ctx context.Context) (*[]Organization, error)
}

// Organization is the struct for the Turso Organization object
type Organization struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Type          string `json:"type"`
	PlanID        string `json:"plan_id"`
	Overages      bool   `json:"overages"`
	BlockedReads  bool   `json:"blocked_reads"`
	BlockedWrites bool   `json:"blocked_writes"`
	PlanTimeline  string `json:"plan_timeline"`
	Memory        int    `json:"memory"`
}

// getOrganizationEndpoint returns the endpoint for the Turso API organization service
func getOrganizationEndpoint(baseURL string) string {
	return fmt.Sprintf("%s/%s", baseURL, organizationEndpoint)
}

// ListOrganizations satisfies the organizationService interface
func (s *OrganizationService) ListOrganizations(ctx context.Context) (*[]Organization, error) {
	endpoint := getOrganizationEndpoint(s.client.cfg.BaseURL)

	resp, err := s.client.DoRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var out []Organization
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("organizations", "listing", resp.StatusCode)
	}

	return &out, nil
}
