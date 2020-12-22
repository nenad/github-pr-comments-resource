package github

import (
	"testing"
)

func TestComment_HasLabel(t *testing.T) {
	tests := []struct {
		name    string
		label   string
		comment Comment
		want    bool
	}{
		{
			name:    "empty comment does not contain empty label",
			comment: Comment{},
			label:   "",
			want:    false,
		},
		{
			name: "comment contains test label",
			comment: Comment{
				Labels: map[string]struct{}{"test": {}},
			},
			label: "test",
			want:  true,
		},
		{
			name: "comment contains other labels but not test",
			comment: Comment{
				Labels: map[string]struct{}{
					"test1": {},
					"test2": {},
				},
			},
			label: "test",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				ID:                tt.comment.ID,
				Body:              tt.comment.Body,
				Labels:            tt.comment.Labels,
				PullRequestNumber: tt.comment.PullRequestNumber,
			}
			if got := c.HasLabel(tt.label); got != tt.want {
				t.Errorf("HasLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}
