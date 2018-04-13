package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
)

func main() {
	client := slack.New("xoxp-2530146065-249234083635-344023717697-59ba56c08afc621ce37e9cdb6b1c7f25")
	params := slack.NewGetFilesParameters()
	params.Channel = "GA4EU44FJ"
	params.Page = 2
	spew.Dump(params)
	files, _, err := client.GetFiles(params)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		log.Println(int64(file.Created), file.Created.Time())
	}
	spew.Dump(files)
}
