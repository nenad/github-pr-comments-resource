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
		quitf("Missing command, should be one of in, out or check")
	}

	ctx := context.Background()

	command := flag.Arg(0)
	var output interface{}

	req := concourse.Request{}
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		quitf("could not decode request: %s", err)
	}
	owner, name, err := req.Source.OwnerAndName()
	if err != nil {
		quitf("could not extract repository info: %s", err)
	}
	client := github.NewClient(ctx, req.Source.AccessToken, owner, name)

	switch command {
	case "in":
		if flag.NArg() < 2 {
			quitf("Missing directory argument")
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
