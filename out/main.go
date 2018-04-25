package main

import (
	"encoding/json"
	"log"
	"os"

	"code.cloudfoundry.org/slack-attachment-resource/shared"
)

type Input struct {
	Source  shared.Source  `json:"source"`
	Version shared.Version `json:"version"`
}

type Output struct {
	Version shared.Version `json:"version"`
}

func main() {
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		log.Fatal("reading input from stdin:", err)
	}

	output := Output{Version: input.Version}
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		log.Fatal("writing output to stdout:", err)
	}
}
