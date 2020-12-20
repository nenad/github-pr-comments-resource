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

	CheckRequest struct {
		Source  Source  `json:"source"`
		Version Version `json:"version"`
	}
)

func (s *Source) OwnerAndName() (string, string, error) {
	parts := strings.Split(s.Repository, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("repository should be in the owner/repository format, got %q", s.Repository)
	}

	return parts[0], parts[1], nil
}
