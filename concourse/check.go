package concourse

import (
	"context"
	"fmt"
	"sort"

	"github.com/nenad/github-pr-comments-resource/github"
)

type (
	CommentSelector interface {
		SelectComments(context.Context, ...github.Filter) ([]github.Comment, error)
	}

	Check struct {
		gh CommentSelector
	}
)

func NewCheck(gh CommentSelector) *Check {
	return &Check{
		gh: gh,
	}
}

func (c *Check) Run(ctx context.Context, request Request) ([]Version, error) {
	var filters []github.Filter
	if request.Version.IDNumber() != 0 {
		filters = append(filters, github.FilterSince(request.Version.IDNumber()))
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

	comments, err := c.gh.SelectComments(ctx, filters...)
	if err != nil {
		return nil, fmt.Errorf("could not fetch comments: %w", err)
	}

	versions := make([]Version, len(comments))
	for idx, com := range comments {
		versions[idx] = VersionFromNumber(com.ID)
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return versions[i].CommentID < versions[j].CommentID
	})

	return versions, nil
}
