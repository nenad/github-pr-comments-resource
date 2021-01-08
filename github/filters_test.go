package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterCommentRegex(t *testing.T) {
	var allComments = []Comment{
		{ID: 1, Body: "test-comment"},
		{ID: 2, Body: "just-a-comment"},
		{ID: 3, Body: "no-comment"},
		{ID: 4, Body: "test-comment"},
	}

	tests := []struct {
		name     string
		regex    string
		comments []Comment
		want     []Comment
	}{
		{
			name:     "comments that match regex",
			regex:    "^test-comment$",
			comments: allComments,
			want: []Comment{
				allComments[0],
				allComments[3],
			},
		},
		{
			name:     "no comments when they do not match regex",
			regex:    "^does-not-match-regex$",
			comments: allComments,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, FilterCommentRegex(tt.regex)(tt.comments))
			},
		)
	}
}

func TestFilterLabels(t *testing.T) {
	var allComments = []Comment{
		{ID: 1, Labels: map[string]struct{}{"test-label": {}}},
		{ID: 2, Labels: map[string]struct{}{"very-wrong-label": {}}},
		{ID: 3, Labels: map[string]struct{}{"test-labelll": {}}},
		{ID: 4, Labels: map[string]struct{}{"test-label": {}}},
	}

	tests := []struct {
		name     string
		labels   []string
		comments []Comment
		want     []Comment
	}{
		{
			name:     "comments that match labels",
			labels:   []string{"test-label"},
			comments: allComments,
			want:     []Comment{allComments[0], allComments[3]},
		},
		{
			name:     "no comments that match no labels",
			labels:   []string{"does-not-exist"},
			comments: allComments,
			want:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, FilterLabels(tt.labels)(tt.comments))
			},
		)
	}
}

func TestFilterLatestPerPR(t *testing.T) {
	var allComments = []Comment{
		{ID: 1, PullRequestNumber: 1},
		{ID: 2, PullRequestNumber: 2},
		{ID: 3, PullRequestNumber: 2},
		{ID: 4, PullRequestNumber: 1},
	}

	tests := []struct {
		name     string
		comments []Comment
		want     []Comment
	}{
		{
			name:     "result should be latest comments of two PRs",
			comments: allComments,
			want:     []Comment{allComments[3], allComments[2]},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, FilterLatestPerPR()(tt.comments))
			},
		)
	}
}

func TestFilterSince(t *testing.T) {
	var allComments = []Comment{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
	}

	tests := []struct {
		name      string
		commentID uint
		comments  []Comment
		want      []Comment
	}{
		{
			name:      "result should only have newer comments than ID 2",
			commentID: 2,
			comments:  allComments,
			want:      []Comment{allComments[2], allComments[3]},
		},
		{
			name:      "no comments should be returned if newest ID is 4",
			commentID: 4,
			comments:  allComments,
			want:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, FilterSince(tt.commentID)(tt.comments))
			},
		)
	}
}
