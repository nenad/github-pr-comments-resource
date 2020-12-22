package concourse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/nenad/github-pr-comments-resource/github"
)

const metaPRNumber = "pr_number"

type (
	In struct {
		gh CommentFetcher
	}

	CommentFetcher interface {
		GetComment(ctx context.Context, commentID uint) (github.Comment, error)
	}
)

func NewIn(gh CommentFetcher) *In {
	return &In{gh}
}

func (i *In) Run(ctx context.Context, req Request, writer io.Writer) (Response, error) {
	c, err := i.gh.GetComment(ctx, req.Version.IDNumber())
	if err != nil {
		return Response{}, fmt.Errorf("could not fetch comment: %w", err)
	}

	var labels []string
	for l := range c.Labels {
		labels = append(labels, l)
	}

	resp := Response{
		Version: req.Version,
		Metadata: []Metadata{
			{Key: metaPRNumber, Value: strconv.FormatUint(uint64(c.PullRequestNumber), 10)},
			{Key: "body", Value: c.Body},
			{Key: "id", Value: strconv.FormatUint(uint64(c.ID), 10)},
			{Key: "labels", Value: strings.Join(labels, ",")},
		},
	}

	if err := json.NewEncoder(writer).Encode(resp); err != nil {
		return Response{}, fmt.Errorf("could not serialize 'in' response: %w", err)
	}

	return resp, nil
}
