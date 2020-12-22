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

type (
	Out struct {
		gh CommentPusher
	}

	CommentPusher interface {
		PutComment(ctx context.Context, body string, prNumber int) (github.Comment, error)
	}
)

func NewOut(gh CommentPusher) *Out {
	return &Out{gh}
}

func (o *Out) Run(ctx context.Context, req Request, reader io.Reader) (Response, error) {
	var inResponse Response
	if err := json.NewDecoder(reader).Decode(&inResponse); err != nil {
		return Response{}, fmt.Errorf("could not decode previous response: %w", err)
	}

	body := req.Params["body"]

	prNum, err := strconv.Atoi(inResponse.GetMetadataField(metaPRNumber))
	if err != nil {
		return Response{}, fmt.Errorf("could not parse PR number from metadata: %w", err)
	}

	c, err := o.gh.PutComment(ctx, body, prNum)
	if err != nil {
		return Response{}, fmt.Errorf("could not post comment to PR: %w", err)
	}

	var labels []string
	for l := range c.Labels {
		labels = append(labels, l)
	}

	outResponse := Response{
		Version: inResponse.Version,
		Metadata: []Metadata{
			{Key: metaPRNumber, Value: strconv.FormatUint(uint64(c.PullRequestNumber), 10)},
			{Key: "body", Value: c.Body},
			{Key: "id", Value: strconv.FormatUint(uint64(c.ID), 10)},
			{Key: "labels", Value: strings.Join(labels, ",")},
		},
	}

	return outResponse, nil
}
