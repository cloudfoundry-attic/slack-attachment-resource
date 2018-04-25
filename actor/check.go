package actor

import (
	"sort"
	"strconv"

	"code.cloudfoundry.org/slack-attachment-resource/shared"
	"github.com/nlopes/slack"
)

//go:generate counterfeiter . CheckAPIClient

type CheckAPIClient interface {
	GetFiles(params slack.GetFilesParameters) ([]slack.File, *slack.Paging, error)
}

func Check(client CheckAPIClient, groupID string, filename string, ts string) ([]shared.Version, error) {
	params := slack.NewGetFilesParameters()
	params.Channel = groupID
	if ts != "" {
		timestamp, err := strconv.Atoi(ts)
		if err != nil {
			return nil, err
		}
		params.TimestampFrom = slack.JSONTime(timestamp)
	}

	var files []slack.File
	for {
		resFiles, _, err := client.GetFiles(params)
		if err != nil {
			return nil, err
		}
		files = append(files, resFiles...)

		if len(resFiles) < 100 {
			break
		}
		params.Page++
	}

	var versions []shared.Version
	for _, file := range files {
		if file.Name == filename {
			versions = append(versions, shared.Version{
				ID:        file.ID,
				Timestamp: strconv.FormatInt(int64(file.Created), 10),
			})
		}
	}
	sort.Slice(versions, func(i, j int) bool { return versions[i].Timestamp < versions[j].Timestamp })
	return versions, nil
}
