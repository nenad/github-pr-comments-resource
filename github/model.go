package github

type (
	// Comment contains a representation of GitHub comment.
	Comment struct {
		// ID refers to the database ID for a particular GitHub comment.
		ID uint
		// Body holds the contents of the comment.
		Body string
		// Labels is a map of label names where the key is the name and the value is empty. It serves the purpose
		// for fast querying if a comment contains a particular label.
		Labels map[string]struct{}
		// PullRequestNumber is the reference to the pull request where this comment was made.
		PullRequestNumber uint
	}
)

// HasLabel checks if the PR where the comment was posted contains a given label.
func (c *Comment) HasLabel(label string) bool {
	_, ok := c.Labels[label]

	return ok
}
