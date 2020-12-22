package concourse

import (
	"fmt"
	"strings"
)

type (
	Source struct {
		AccessToken string   `json:"access_token"`
		Repository  string   `json:"repository"`
		Comments    string   `json:"comments"`
		Labels      []string `json:"labels"`
		LatestPerPR bool     `json:"latest_per_pr"`
	}

	Metadata struct {
		Key   string `json:"name"`
		Value string `json:"value"`
	}

	// Request is a generic resource request for check, in and out tasks. Fields are populated per task:
	// check:	Source and Version
	// in:		Source, Version and Params
	// out:		Source and Params
	Request struct {
		Source  Source            `json:"source"`
		Version Version           `json:"version"`
		Params  map[string]string `json:"params"`
	}

	// Response is a generic resource response for in and out tasks.
	Response struct {
		Version  Version    `json:"version"`
		Metadata []Metadata `json:"metadata"`
	}
)

func (s *Source) OwnerAndName() (string, string, error) {
	parts := strings.Split(s.Repository, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("repository should be in the owner/repository format, got %q", s.Repository)
	}

	return parts[0], parts[1], nil
}

func (r *Response) GetMetadataField(name string) string {
	for _, m := range r.Metadata {
		if m.Key == name {
			return m.Value
		}
	}

	return ""
}
