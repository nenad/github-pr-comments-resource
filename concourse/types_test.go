package concourse_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nenad/github-pr-comments-resource/concourse"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name    string
		payload io.Reader
		want    concourse.Request
		wantErr bool
	}{
		{
			name:    "nil data results in error",
			payload: nil,
			want:    concourse.Request{},
			wantErr: true,
		},
		{
			name:    "bad data results in error",
			payload: bytes.NewBufferString(`{"just": "baddata"}`),
			want:    concourse.Request{},
			wantErr: true,
		},
		{
			name: "valid data results in correct full structure",
			payload: bytes.NewBufferString(
				`{
				"source": {
					"repository": "testowner/testrepository",
					"access_token": "testtoken",
					"comments": "hello world",
					"labels": ["test", "label"],
					"latest_per_pr": true
				},
				"version": {
					"comment_id": "1234567"
				},
				"params": {
					"key": "value"
				}
			}`,
			),
			want: concourse.Request{
				Source: concourse.Source{
					AccessToken: "testtoken",
					Repository:  "testowner/testrepository",
					Comments:    "hello world",
					Labels:      []string{"test", "label"},
					LatestPerPR: true,
				},
				Version: concourse.Version{CommentID: "1234567"},
				Params:  map[string]string{"key": "value"},
			},
			wantErr: false,
		},
		{
			name: "invalid repository results in an error",
			payload: bytes.NewBufferString(
				`{"source":{"repository":"badrepo","access_token":"testtoken"}}`,
			),
			want:    concourse.Request{},
			wantErr: true,
		},
		{
			name: "bad characters in repository results in an error",
			payload: bytes.NewBufferString(
				`{"source":{"repository":"test/123&#&@","access_token":"testtoken"}}`,
			),
			want:    concourse.Request{},
			wantErr: true,
		},
		{
			name: "empty access token results in an error",
			payload: bytes.NewBufferString(
				`{"source":{"repository":"owner/repo"}}`,
			),
			want:    concourse.Request{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := concourse.NewRequest(tt.payload)
				if (err != nil) != tt.wantErr {
					t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				assert.Equal(t, got, tt.want)
			},
		)
	}
}

func TestRequest_RepositoryName(t *testing.T) {
	data := bytes.NewBufferString(`{"source":{"repository":"owner/reponame","access_token":"testtoken"}}`)
	r, _ := concourse.NewRequest(data)

	assert.Equal(t, "reponame", r.RepositoryName())
	assert.Equal(t, "owner", r.RepositoryOwner())
}
