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

func TestListGroups(t *testing.T) {
	groupService := newMockGroupService()

	resp, err := groupService.ListGroups(context.Background())
	require.NoError(t, err)
	assert.Len(t, resp.Groups, 1)
}
func TestGetGroup(t *testing.T) {
	groupService := newMockGroupService()

	resp, err := groupService.GetGroup(context.Background(), "meow")
	require.NoError(t, err)
	assert.Equal(t, resp.Group.Name, "meow")
}

func TestDeleteGroup(t *testing.T) {
	groupService := newMockGroupService()

	resp, err := groupService.DeleteGroup(context.Background(), "meow")
	require.NoError(t, err)
	assert.Equal(t, resp.Group.Name, "woof")
	assert.True(t, resp.Group.Archived)
}

func TestCreateGroup(t *testing.T) {
	body := `{"group":{"archived":false,"locations":["lhr","ams","bos"],"name":"meow","primary":"lhr","uuid":"0a28102d-6906-11ee-8553-eaa7715aeaf2","version":"v0.23.7"}}`
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
	groupService := GroupService{client: client}
	req := CreateGroupRequest{
		Name:     "meow",
		Location: "ams",
	}

	resp, err := groupService.CreateGroup(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, resp.Group.Name, "meow")

	// test error
	req = CreateGroupRequest{}

	resp, err = groupService.CreateGroup(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAddLocation(t *testing.T) {
	body := `{"group":{"archived":false,"locations":["lhr","ams","bos", "den"],"name":"meow","primary":"lhr","uuid":"0a28102d-6906-11ee-8553-eaa7715aeaf2","version":"v0.23.7"}}`
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
	groupService := GroupService{client: client}
	req := GroupLocationRequest{
		GroupName: "meow",
		Location:  "den",
	}

	resp, err := groupService.AddLocation(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, resp.Group.Name, "meow")

	// test error, missing location
	req = GroupLocationRequest{
		GroupName: "meow",
	}

	resp, err = groupService.AddLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// test error, missing group name
	req = GroupLocationRequest{
		Location: "den",
	}

	resp, err = groupService.AddLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestRemoveLocation(t *testing.T) {
	body := `{"group":{"archived":false,"locations":["lhr","ams","bos"] ,"name":"meow","primary":"lhr","uuid":"0a28102d-6906-11ee-8553-eaa7715aeaf2","version":"v0.23.7"}}`
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
	groupService := GroupService{client: client}
	req := GroupLocationRequest{
		GroupName: "meow",
		Location:  "den",
	}

	resp, err := groupService.RemoveLocation(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, resp.Group.Name, "meow")

	// test error, missing location
	req = GroupLocationRequest{
		GroupName: "meow",
	}

	resp, err = groupService.RemoveLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// test error, missing group name
	req = GroupLocationRequest{
		Location: "den",
	}

	resp, err = groupService.RemoveLocation(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
func TestValidateGroupCreateRequest(t *testing.T) {
	tests := []struct {
		name    string
		request CreateGroupRequest
		wantErr error
	}{
		{
			name: "Valid request",
			request: CreateGroupRequest{
				Name:     "meow",
				Location: "ams",
			},
			wantErr: nil,
		},
		{
			name: "missing name",
			request: CreateGroupRequest{
				Name:     "",
				Location: "ams",
			},
			wantErr: &MissingRequiredFieldError{RequiredField: "name"},
		},
		{
			name: "invalid name",
			request: CreateGroupRequest{
				Name:     "my group",
				Location: "ams",
			},
			wantErr: &InvalidFieldError{Field: "name", Message: "spaces are not allowed"},
		},
		{
			name: "missing location",
			request: CreateGroupRequest{
				Name:     "the-best",
				Location: "",
			},
			wantErr: &MissingRequiredFieldError{RequiredField: "location"},
		},
		{
			name: "invalid location",
			request: CreateGroupRequest{
				Name:     "the-best",
				Location: "us",
			},
			wantErr: &InvalidFieldError{Field: "location", Message: "must be 3 characters"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGroupCreateRequest(tt.request)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr.Error())

				return
			}

			require.NoError(t, err)
		})
	}
}
