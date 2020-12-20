package github

type (
	pullRequestCommentQuery struct {
		Repository struct {
			PullRequests struct {
				Edges []struct {
					Node struct {
						Number uint `json:"number"`
						Labels struct {
							Edges []struct {
								Node struct {
									Name string `graphql:"name"`
								} `graphql:"node"`
							} `graphql:"edges"`
						} `graphql:"labels(first: $labelsFirst)"`
						Comments struct {
							Edges []struct {
								Cursor string `graphql:"cursor"`
								Node   struct {
									Body string `graphql:"body"`
									ID   uint   `graphql:"databaseId"`
								} `graphql:"node"`
							} `graphql:"edges"`
						} `graphql:"comments(last: $commentsLast)"`
					} `graphql:"node"`
				} `graphql:"edges"`
			} `graphql:"pullRequests(last: $pullRequestsLast, states: $states)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
)

func (p *pullRequestCommentQuery) toComments() []Comment {
	var comments []Comment
	for _, pr := range p.Repository.PullRequests.Edges {
		for _, com := range pr.Node.Comments.Edges {
			labels := make(map[string]struct{}, len(pr.Node.Labels.Edges))
			for _, lab := range pr.Node.Labels.Edges {
				labels[lab.Node.Name] = struct{}{}
			}

			comments = append(comments, Comment{
				ID:                com.Node.ID,
				Body:              com.Node.Body,
				Labels:            labels,
				PullRequestNumber: pr.Node.Number,
			})
		}
	}

	return comments
}
