package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/nenad/github-pr-comments-resource/concourse"
	"github.com/nenad/github-pr-comments-resource/github"
)

func runCheck(ctx context.Context, request concourse.CheckRequest) ([]concourse.Version, error) {
	owner, repo, err := request.Source.OwnerAndName()
	if err != nil {
		return nil, fmt.Errorf("could not extract repository information: %w", err)
	}

	c := github.NewClient(ctx, request.Source.AccessToken, owner, repo)

	var filters []github.Filter
	if request.Version.CommentID != 0 {
		filters = append(filters, github.FilterSince(request.Version.CommentID))
	}

	if request.Source.Comments != "" {
		filters = append(filters, github.FilterCommentRegex(request.Source.Comments))
	}

	if len(request.Source.Labels) != 0 {
		filters = append(filters, github.FilterLabels(request.Source.Labels))
	}

	if request.Source.LatestPerPR {
		filters = append(filters, github.FilterLatestPerPR())
	}

	comments, err := c.SelectComments(ctx, filters...)
	if err != nil {
		return nil, fmt.Errorf("could not fetch comments: %w", err)
	}

	versions := make([]concourse.Version, len(comments))
	for idx, com := range comments {
		versions[idx] = concourse.Version{
			CommentID: com.ID,
		}
	}

	// Should sorting be done in client functions?
	sort.SliceStable(versions, func(i, j int) bool {
		return versions[i].CommentID < versions[j].CommentID
	})

	return versions, nil
}
