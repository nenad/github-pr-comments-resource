package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/nenad/github-pr-comments-resource/concourse"
)

// TODO extract parameters and move to function
func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		quit("Missing command, should be one of in, out or check")
	}

	ctx := context.Background()

	command := flag.Arg(0)
	var output interface{}
	var err error
	switch command {
	case "in":
		fmt.Printf("Yay")
	case "out":
		fmt.Printf("Yay")
	case "check":
		req := concourse.CheckRequest{}
		if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
			quit("could not decode check request: %s", err.Error())
		}
		output, err = runCheck(ctx, req)
	default:
		quit("Invalid command, should be one of in, out or check, got %s", command)
	}

	if err != nil {
		quit(err.Error())
	}

	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		quit(err.Error())
	}
}

func quit(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}
