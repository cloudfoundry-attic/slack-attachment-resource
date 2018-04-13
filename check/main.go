package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nlopes/slack"

	"code.cloudfoundry.org/slack-attachment-resource/actor"
	"code.cloudfoundry.org/slack-attachment-resource/shared"
)

type Input struct {
	Source  shared.Source  `json:"source"`
	Version shared.Version `json:"version"`
}

type Output []shared.Version

func main() {
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		log.Fatal("reading input from stdin:", err)
	}

	client := slack.New(input.Source.Token)
	versions, err := actor.Check(client, input.Source.GroupID, input.Source.Filename, input.Version.Timestamp)
	if err != nil {
		log.Fatal("running check with slack client:", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(versions); err != nil {
		log.Fatal("writing output to stdout:", err)
	}
}
