package main

import (
	"encoding/json"
	"log"
	"os"

	"code.cloudfoundry.org/slack-attachment-resource/actor"
	"code.cloudfoundry.org/slack-attachment-resource/shared"
	"github.com/nlopes/slack"
)

type Input struct {
	Source  shared.Source  `json:"source"`
	Version shared.Version `json:"version"`
}

type Output struct {
	Version shared.Version `json:"version"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must pass a output directory")
	}

	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		log.Fatal("reading input from stdin:", err)
	}

	client := slack.New(input.Source.Token)
	err := actor.In(client, input.Source.Token, input.Version, os.Args[1])
	if err != nil {
		log.Fatal("running in with slack client:", err)
	}

	output := Output{Version: input.Version}
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		log.Fatal("writing output to stdout:", err)
	}
}
