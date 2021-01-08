package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nenad/github-pr-comments-resource/concourse"
	"github.com/nenad/github-pr-comments-resource/github"
)

// TODO extract parameters and move to function.
func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		quitf("missing command, should be one of in, out or check")
	}


	req, err := concourse.NewRequest(os.Stdin)
	if err != nil {
		quitf("could not create a new request: %s", err)
	}

	ctx := context.Background()
	client := github.NewClient(ctx, req.Source.AccessToken, req.RepositoryOwner(), req.RepositoryName())

	var output interface{}
	command := flag.Arg(0)
	switch command {
	case "in":
		if flag.NArg() < 2 {
			quitf("missing directory argument")
		}

		// TODO Configurable file name?
		dir := flag.Arg(1)
		metafile, err := os.Create(filepath.Join(dir, "metadata.json"))
		if err != nil {
			quitf("could not write 'in' metadata: %s", err)
		}

		in := concourse.NewIn(client)
		output, err = in.Run(ctx, req, metafile)
		if err != nil {
			quitf("in failed: %s", err)
		}
	case "out":
		if flag.NArg() < 2 {
			quitf("Missing directory argument")
		}

		// TODO Configurable file name?
		dir := flag.Arg(1)
		metafile, err := os.Open(filepath.Join(dir, "metadata.json"))
		if err != nil {
			quitf("could not read 'out' metadata: %s", err)
		}

		out := concourse.NewOut(client)
		output, err = out.Run(ctx, req, metafile)
		if err != nil {
			quitf("out failed: %s", err)
		}
	case "check":
		check := concourse.NewCheck(client)
		output, err = check.Run(ctx, req)
		if err != nil {
			quitf("check failed: %s", err)
		}
	default:
		quitf("Invalid command, should be one of in, out or check, got %s", command)
	}

	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		quitf("could not encode response: %s", err)
	}
}

func quitf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}
