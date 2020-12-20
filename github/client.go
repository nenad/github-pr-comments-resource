package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type (
	Client struct {
		owner          string
		repositoryName string
		gh             *githubv4.Client
	}

	Filter func([]Comment) []Comment
)

func NewClient(ctx context.Context, accessToken, owner, repositoryName string) *Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	httpClient := oauth2.NewClient(ctx, src)

	client := githubv4.NewClient(httpClient)
	return &Client{gh: client, owner: owner, repositoryName: repositoryName}
}

func (c *Client) GetComments(ctx context.Context) ([]Comment, error) {
	var queryResult pullRequestCommentQuery
	err := c.gh.Query(ctx, &queryResult, map[string]interface{}{
		"owner":            githubv4.String(c.owner),
		"name":             githubv4.String(c.repositoryName),
		"pullRequestsLast": githubv4.Int(30), // TODO Configurable value
		"commentsLast":     githubv4.Int(10), // TODO Configurable value
		"labelsFirst":      githubv4.Int(100),
		"states":           []githubv4.PullRequestState{githubv4.PullRequestStateOpen},
	})
	if err != nil {
		return nil, fmt.Errorf("error from GitHub: %w", err)
	}

	return queryResult.toComments(), nil
}

func (c *Client) SelectComments(ctx context.Context, filters ...Filter) ([]Comment, error) {
	comments, err := c.GetComments(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch all comments: %w", err)
	}

	for _, f := range filters {
		comments = f(comments)
	}

	return comments, nil
}
