package concourse

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	repositoryRegex = regexp.MustCompile(`^[a-zA-Z0-9-_.]+/[a-zA-Z0-9-_.]+$`)
)

type (
	// Source defines the resource configuration.
	Source struct {
		AccessToken string   `json:"access_token"`
		Repository  string   `json:"repository"`
		Comments    string   `json:"comments"`
		Labels      []string `json:"labels"`
		LatestPerPR bool     `json:"latest_per_pr"`
	}

	// Metadata is a metadata entry.
	// Example in: https://concourse-ci.org/implementing-resource-types.html#resource-in
	Metadata struct {
		Key   string `json:"name"`
		Value string `json:"value"`
	}

	// Request is a generic resource request for check, in and out tasks. Fields are populated per task:
	// check:	Source and Version
	// in:		Source, Version and Params
	// out:		Source and Params
	Request struct {
		Source          Source            `json:"source"`
		Version         Version           `json:"version"`
		Params          map[string]string `json:"params"`
	}

	// Response is a generic resource response for in and out tasks.
	Response struct {
		Version  Version    `json:"version"`
		Metadata []Metadata `json:"metadata"`
	}
)

// NewRequest validates and creates a request, or returns an error
// if validations fail.
func NewRequest(reader io.Reader) (Request, error) {
	req := Request{}
	if reader == nil {
		return Request{}, fmt.Errorf("no reader was provided")
	}

	if err := json.NewDecoder(reader).Decode(&req); err != nil {
		return Request{}, fmt.Errorf("could not decode request: %w", err)
	}

	if !repositoryRegex.MatchString(req.Source.Repository) {
		return Request{}, fmt.Errorf("repository must be in the format <repositoryOwner>/<repositoryName>")
	}

	if req.Source.AccessToken == "" {
		return Request{}, fmt.Errorf("access_token must be set")
	}

	return req, nil
}

// RepositoryOwner returns the owner of a repository. For example,
// returns "testowner" from the source repository "testowner/testrepo".
func (r *Request) RepositoryOwner() string {
	return strings.Split(r.Source.Repository, "/")[0]
}

// RepositoryName returns the name of a repository. For example,
// returns "testrepo" from the source repository "testowner/testrepo".
func (r *Request) RepositoryName() string {
	return strings.Split(r.Source.Repository, "/")[1]
}

// GetMetadataField looks up the metadata of a response.
func (r *Response) GetMetadataField(name string) string {
	for _, m := range r.Metadata {
		if m.Key == name {
			return m.Value
		}
	}

	return ""
}
