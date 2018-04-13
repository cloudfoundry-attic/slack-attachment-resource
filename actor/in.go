package actor

import (
	"fmt"
	"net/http"

	"code.cloudfoundry.org/slack-attachment-resource/shared"
	"github.com/nlopes/slack"
)

//go:generate counterfeiter . InAPIClient

type InAPIClient interface {
	GetFileInfo(fileID string, count int, page int) (*slack.File, []slack.Comment, *slack.Paging, error)
}

func In(client InAPIClient, authorizationToken string, version shared.Version, outputDirectory string) error {
	file, _, _, err := client.GetFileInfo(version.ID, 0, 0)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", file.URLPrivateDownload, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authorizationToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	_ = res

	if res.StatusCode >= 400 {
		return HTTPError{Message: res.Status, URL: file.URLPrivateDownload}
	}

	return nil
}

type HTTPError struct {
	Message string
	URL     string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("http error retrieving file from %s: %s", e.URL, e.Message)
}
