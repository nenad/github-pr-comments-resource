package github

type (
	// Comment contains a representation of GitHub comment.
	Comment struct {
		// Body holds the contents of the comment.
		Body string
		// ID refers to the database ID for a particular GitHub comment.
		ID uint
		// Labels is a map of label names where the key is the name and the value is empty. It serves the purpose
		// for fast querying if a comment contains a particular label.
		Labels map[string]struct{}
		// PullRequestNumber is the reference to the pull request where this comment was made.
		PullRequestNumber uint
	}
)

func (c *Comment) HasLabel(label string) bool {
	_, ok := c.Labels[label]
	return ok
}
