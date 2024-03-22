package turso

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	groupEndpoint    = "v1/organizations/%s/groups"
	locationEndpoint = groupEndpoint + "%s/locations/%s"
)

// GroupService is the interface for the Turso API group endpoint
type GroupService service

type groupService interface {
	// ListGroups lists all groups in the organization
	ListGroups(ctx context.Context) (*ListGroupResponse, error)
	// CreateGroup creates a new group in the organization
	CreateGroup(ctx context.Context, req CreateGroupRequest) (*CreateGroupResponse, error)
	// GetGroup gets a group by name
	GetGroup(ctx context.Context, groupName string) (*GetGroupResponse, error)
	// DeleteGroup deletes a group by name
	DeleteGroup(ctx context.Context, groupName string) (*DeleteGroupResponse, error)
	// AddLocation adds a location to a group
	AddLocation(ctx context.Context, eq GroupLocationRequest) (*GroupLocationResponse, error)
	// RemoveLocation removes a location from a group
	RemoveLocation(ctx context.Context, req GroupLocationRequest) (*GroupLocationResponse, error)
}

// Group is the struct for the Turso API group service
type Group struct {
	Archived  bool     `json:"archived"`
	Locations []string `json:"locations"`
	Name      string   `json:"name"`
	Primary   string   `json:"primary"`
	UUID      string   `json:"uuid"`
	Version   string   `json:"version"`
}

// ListGroupResponse is the struct for the Turso API group list response
type ListGroupResponse struct {
	Groups []Group `json:"groups"`
}

// GetGroupResponse is the struct for the Turso API group get response
type GetGroupResponse struct {
	Group Group `json:"group"`
}

// CreateGroupResponse is the struct for the Turso API group create response
type CreateGroupResponse struct {
	Group Group `json:"group"`
}

// GroupLocationRequest is the struct for the Turso API to add or remove a location to a group request
type GroupLocationRequest struct {
	// GroupName is the name of the group to add the location
	GroupName string
	// Location is the location to add to the group
	Location string
}

// GroupLocationResponse is the struct for the Turso API to add or remove location to group response
type GroupLocationResponse struct {
	Group Group `json:"group"`
}

// DeleteGroupResponse is the struct for the Turso API group delete response
type DeleteGroupResponse struct {
	Group Group `json:"group"`
}

// CreateGroupRequest is the struct for the Turso API group create request
type CreateGroupRequest struct {
	Extensions string `json:"extensions"`
	Location   string `json:"location"`
	Name       string `json:"name"`
}

// getGroupEndpoint returns the endpoint for the Turso API group service
func getGroupEndpoint(baseURL, orgName string) string {
	groupEndpoint := fmt.Sprintf(groupEndpoint, orgName)
	return fmt.Sprintf("%s/%s", baseURL, groupEndpoint)
}

// getGroupLocationsEndpoint returns the endpoint for the Turso API group locations service
func getGroupLocationsEndpoint(baseURL, orgName, groupName, locationName string) string {
	locEndpoint := fmt.Sprintf(locationEndpoint, orgName, groupName, locationName)
	return fmt.Sprintf("%s/%s", baseURL, locEndpoint)
}

// ListGroups satisfies the groupService interface
func (s *GroupService) ListGroups(ctx context.Context) (*ListGroupResponse, error) {
	endpoint := getGroupEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName)

	resp, err := s.client.DoRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var out ListGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("groups", "listing", resp.StatusCode)
	}

	return &out, nil
}

// CreateGroup satisfies the groupService interface
func (s *GroupService) CreateGroup(ctx context.Context, group CreateGroupRequest) (*CreateGroupResponse, error) {
	// Validate the request
	if err := validateGroupCreateRequest(group); err != nil {
		return nil, err
	}

	// Create the group
	endpoint := getGroupEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName)

	resp, err := s.client.DoRequest(ctx, http.MethodPost, endpoint, group)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode the response
	var out CreateGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("group", "creating", resp.StatusCode)
	}

	return &out, nil
}

// GetGroup satisfies the groupService interface
func (s *GroupService) GetGroup(ctx context.Context, groupName string) (*GetGroupResponse, error) {
	// get endpoint and append the group name
	endpoint := getGroupEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName)
	endpoint = fmt.Sprintf("%s/%s", endpoint, groupName)

	resp, err := s.client.DoRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var out *GetGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("group", "getting", resp.StatusCode)
	}

	return out, nil
}

// DeleteGroup satisfies the groupService interface
func (s *GroupService) DeleteGroup(ctx context.Context, groupName string) (*DeleteGroupResponse, error) {
	// Create the group
	endpoint := getGroupEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName)
	endpoint = fmt.Sprintf("%s/%s", endpoint, groupName)

	resp, err := s.client.DoRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode the response
	var out DeleteGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("group", "deleting", resp.StatusCode)
	}

	return &out, nil
}

// AddLocation satisfies the groupService interface
func (s *GroupService) AddLocation(ctx context.Context, req GroupLocationRequest) (*GroupLocationResponse, error) {
	if err := validateLocationRequest(req); err != nil {
		return nil, err
	}

	endpoint := getGroupLocationsEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName, req.GroupName, req.Location)

	resp, err := s.client.DoRequest(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode the response
	var out GroupLocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("group", "creating", resp.StatusCode)
	}

	return &out, nil
}

// RemoveLocation satisfies the groupService interface
func (s *GroupService) RemoveLocation(ctx context.Context, req GroupLocationRequest) (*GroupLocationResponse, error) {
	if err := validateLocationRequest(req); err != nil {
		return nil, err
	}

	endpoint := getGroupLocationsEndpoint(s.client.cfg.BaseURL, s.client.cfg.OrgName, req.GroupName, req.Location)

	resp, err := s.client.DoRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode the response
	var out GroupLocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newBadRequestError("group", "creating", resp.StatusCode)
	}

	return &out, nil
}

// validateGroupCreateRequest validates the group create request
func validateGroupCreateRequest(req CreateGroupRequest) error {
	if err := validateGroupName(req.Name); err != nil {
		return err
	}

	if err := validateLocation(req.Location); err != nil {
		return err
	}

	return nil
}

// validateGroupName validates the group name
func validateGroupName(groupName string) error {
	if groupName == "" {
		return newMissingRequiredFieldError("name")
	}

	if strings.Contains(groupName, " ") {
		return newInvalidFieldError("name", "spaces are not allowed")
	}

	return nil
}

// validateGroupName validates the group name
func validateLocation(location string) error {
	if location == "" {
		return newMissingRequiredFieldError("location")
	}

	// all Turso locations are 3 characters
	if len(location) != 3 { //nolint:gomnd
		return newInvalidFieldError("location", "must be 3 characters")
	}

	return nil
}

// validateLocationRequest validates the group location request fields
func validateLocationRequest(req GroupLocationRequest) error {
	if req.GroupName == "" {
		return newMissingRequiredFieldError("name")
	}

	if req.Location == "" {
		return newMissingRequiredFieldError("location")
	}

	return nil
}
