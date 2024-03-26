package turso

import (
	"context"
	"net/http"
)

// MockHTTPRequestDoer implements the standard http.Client interface.
type MockHTTPRequestDoer struct {
	Response *http.Response
	Error    error
}

// Do implements the standard http.Client interface for MockHTTPRequestDoer
func (md *MockHTTPRequestDoer) Do(req *http.Request) (*http.Response, error) {
	return md.Response, md.Error
}

// NewMockClient creates a new client for interacting with the Turso API to mock ok requests
// this can be used to test the client without hitting the actual API an expect an 200 OK response.
func NewMockClient() *Client {
	c := &Client{}
	c.Group = newMockGroupService()
	c.Database = newMockDatabaseService()
	c.Organization = newMockOrganizationService()
	c.DatabaseTokens = newMockDatabaseTokenService()

	return c
}

type MockGroupService struct {
	ListGroupResponse     *ListGroupResponse
	CreateGroupResponse   *CreateGroupResponse
	GetGroupResponse      *GetGroupResponse
	DeleteGroupResponse   *DeleteGroupResponse
	GroupLocationResponse *GroupLocationResponse
	Error                 error
}

type MockDatabaseService struct {
	ListDatabaseResponse   *ListDatabaseResponse
	CreateDatabaseResponse *CreateDatabaseResponse
	GetDatabaseResponse    *GetDatabaseResponse
	DeleteDatabaseResponse *DeleteDatabaseResponse
	Error                  error
}

type MockDatabaseTokensService struct {
	CreateDatabaseTokenResponse *CreateDatabaseTokenResponse
	Error                       error
}

type MockOrganizationService struct {
	ListOrganizationsResponse *[]Organization
	Error                     error
}

func newMockGroupService() groupService {
	return &MockGroupService{
		ListGroupResponse: &ListGroupResponse{
			Groups: []Group{
				{
					Archived:  false,
					Locations: []string{"lhr", "ams", "bos"},
					Name:      "meow",
					Primary:   "lhr",
					UUID:      "0a28102d-6906-11ee-8553-eaa7715aeaf2",
					Version:   "v0.23.7",
				},
			},
		},
		CreateGroupResponse: &CreateGroupResponse{
			Group: Group{
				Archived:  false,
				Locations: []string{"lhr", "ams", "bos"},
				Name:      "meow",
				Primary:   "lhr",
				UUID:      "0a28102d-6906-11ee-8553-eaa7715aeaf2",
				Version:   "v0.23.7",
			},
		},
		GetGroupResponse: &GetGroupResponse{
			Group: Group{
				Archived:  false,
				Locations: []string{"lhr", "ams", "bos"},
				Name:      "meow",
				Primary:   "lhr",
				UUID:      "0a28102d-6906-11ee-8553-eaa7715aeaf2",
				Version:   "v0.23.7",
			},
		},
		DeleteGroupResponse: &DeleteGroupResponse{
			Group: Group{
				Archived:  true,
				Locations: []string{"lhr", "ams", "bos"},
				Name:      "woof",
				Primary:   "lhr",
				UUID:      "0a28102d-6906-11ee-8553-eaa7715aeaf2",
				Version:   "v0.23.7",
			},
		},
		GroupLocationResponse: &GroupLocationResponse{
			Group: Group{
				Archived:  false,
				Locations: []string{"lhr", "ams", "bos"},
				Name:      "meow",
				Primary:   "lhr",
				UUID:      "0a28102d-6906-11ee-8553-eaa7715aeaf2",
				Version:   "v0.23.7",
			},
		},
		Error: nil,
	}
}

func newMockDatabaseService() databaseService {
	return &MockDatabaseService{
		ListDatabaseResponse: &ListDatabaseResponse{
			Databases: []*Database{
				{
					Name:       "my-db",
					Hostname:   "[databaseName]-[organizationName].turso.io",
					DatabaseID: "0eb771dd-6906-11ee-8553-eaa7715aeaf2",
				},
			},
		},
		CreateDatabaseResponse: &CreateDatabaseResponse{
			CreateDatabase{
				Name:       "my-db",
				Hostname:   "[databaseName]-[organizationName].turso.io",
				DatabaseID: "0eb771dd-6906-11ee-8553-eaa7715aeaf2",
			},
		},
		GetDatabaseResponse: &GetDatabaseResponse{
			Database: &Database{
				Name:       "my-db",
				Hostname:   "[databaseName]-[organizationName].turso.io",
				DatabaseID: "0eb771dd-6906-11ee-8553-eaa7715aeaf2",
			},
		},
		DeleteDatabaseResponse: &DeleteDatabaseResponse{
			Database: "my-db",
		},
		Error: nil,
	}
}

func newMockOrganizationService() organizationService {
	return &MockOrganizationService{
		ListOrganizationsResponse: &[]Organization{
			{
				Name: "meow",
				Slug: "meow",
			},
		},
		Error: nil,
	}
}

func newMockDatabaseTokenService() databaseTokensService {
	return &MockDatabaseTokensService{
		CreateDatabaseTokenResponse: &CreateDatabaseTokenResponse{
			JWT: "jwt-token",
		},
		Error: nil,
	}
}

func (mg *MockGroupService) ListGroups(ctx context.Context) (*ListGroupResponse, error) {
	return mg.ListGroupResponse, mg.Error
}

func (mg *MockGroupService) CreateGroup(ctx context.Context, req CreateGroupRequest) (*CreateGroupResponse, error) {
	return mg.CreateGroupResponse, mg.Error
}

func (mg *MockGroupService) GetGroup(ctx context.Context, groupName string) (*GetGroupResponse, error) {
	return mg.GetGroupResponse, mg.Error
}

func (mg *MockGroupService) DeleteGroup(ctx context.Context, groupName string) (*DeleteGroupResponse, error) {
	return mg.DeleteGroupResponse, mg.Error
}

func (mg *MockGroupService) AddLocation(ctx context.Context, eq GroupLocationRequest) (*GroupLocationResponse, error) {
	return mg.GroupLocationResponse, mg.Error
}

func (mg *MockGroupService) RemoveLocation(ctx context.Context, req GroupLocationRequest) (*GroupLocationResponse, error) {
	return mg.GroupLocationResponse, mg.Error
}

func (md *MockDatabaseService) ListDatabases(ctx context.Context) (*ListDatabaseResponse, error) {
	return md.ListDatabaseResponse, md.Error
}

func (md *MockDatabaseService) CreateDatabase(ctx context.Context, req CreateDatabaseRequest) (*CreateDatabaseResponse, error) {
	return md.CreateDatabaseResponse, md.Error
}

func (md *MockDatabaseService) GetDatabase(ctx context.Context, dbName string) (*GetDatabaseResponse, error) {
	return md.GetDatabaseResponse, md.Error
}

func (md *MockDatabaseService) DeleteDatabase(ctx context.Context, dbName string) (*DeleteDatabaseResponse, error) {
	return md.DeleteDatabaseResponse, md.Error
}

func (mo *MockOrganizationService) ListOrganizations(ctx context.Context) (*[]Organization, error) {
	return mo.ListOrganizationsResponse, mo.Error
}

func (md *MockDatabaseTokensService) CreateDatabaseToken(ctx context.Context, req CreateDatabaseTokenRequest) (*CreateDatabaseTokenResponse, error) {
	return md.CreateDatabaseTokenResponse, md.Error
}
