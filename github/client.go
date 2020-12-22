package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v33/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type (
	Client struct {
		owner          string
		repositoryName string
		ghv3           *github.Client
		ghv4           *githubv4.Client
	}

	Filter func([]Comment) []Comment
)

func NewClient(ctx context.Context, accessToken, owner, repositoryName string) *Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	authClient := oauth2.NewClient(ctx, src)

	v3client := github.NewClient(authClient)
	v4client := githubv4.NewClient(authClient)

	return &Client{ghv4: v4client, ghv3: v3client, owner: owner, repositoryName: repositoryName}
}

func (c *Client) GetComments(ctx context.Context) ([]Comment, error) {
	var queryResult pullRequestCommentQuery
	err := c.ghv4.Query(ctx, &queryResult, map[string]interface{}{
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

func (c *Client) GetComment(ctx context.Context, commentID uint) (Comment, error) {
	ic, _, err := c.ghv3.Issues.GetComment(ctx, c.owner, c.repositoryName, int64(commentID))
	if err != nil {
		return Comment{}, fmt.Errorf("could not get comment on issue: %w", err)
	}

	prNumber := extractPRNumber(ic.GetIssueURL())
	pr, _, err := c.ghv3.PullRequests.Get(ctx, c.owner, c.repositoryName, prNumber)
	if err != nil {
		return Comment{}, fmt.Errorf("could not get pull request information: %w", err)
	}

	labels := make(map[string]struct{})
	for _, l := range pr.Labels {
		labels[l.GetName()] = struct{}{}
	}

	return Comment{
		Body:              ic.GetBody(),
		ID:                uint(ic.GetID()),
		PullRequestNumber: uint(pr.GetNumber()),
		Labels:            labels,
	}, nil
}

func (c *Client) PutComment(ctx context.Context, body string, prNumber int) (Comment, error) {
	comment := &github.IssueComment{
		Body: &body,
	}
	ic, _, err := c.ghv3.Issues.CreateComment(ctx, c.owner, c.repositoryName, prNumber, comment)
	if err != nil {
		return Comment{}, fmt.Errorf("could not get comment on issue: %w", err)
	}

	pr, _, err := c.ghv3.PullRequests.Get(ctx, c.owner, c.repositoryName, prNumber)
	if err != nil {
		return Comment{}, fmt.Errorf("could not get pull request information: %w", err)
	}

	labels := make(map[string]struct{})
	for _, l := range pr.Labels {
		labels[l.GetName()] = struct{}{}
	}

	return Comment{
		Body:              ic.GetBody(),
		ID:                uint(ic.GetID()),
		PullRequestNumber: uint(pr.GetNumber()),
		Labels:            labels,
	}, nil
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

func extractPRNumber(url string) int {
	parts := strings.Split(url, "/")
	idStr := parts[len(parts)-1]
	id, _ := strconv.Atoi(idStr)

	return id
}
