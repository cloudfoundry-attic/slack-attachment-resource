package actor_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	. "code.cloudfoundry.org/slack-attachment-resource/actor"
	"code.cloudfoundry.org/slack-attachment-resource/actor/actorfakes"
	"code.cloudfoundry.org/slack-attachment-resource/shared"
	"github.com/nlopes/slack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

var _ = FDescribe("In", func() {
	var (
		fakeClient *actorfakes.FakeInAPIClient
		token      string
		version    shared.Version
		outputDir  string

		executeErr error
	)

	BeforeEach(func() {
		fakeClient = new(actorfakes.FakeInAPIClient)
		token = "don't token me bro"
		version = shared.Version{
			ID:        "12345",
			Timestamp: "6789",
		}

		var err error
		outputDir, err = ioutil.TempDir("", "in-test")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(outputDir)).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		executeErr = In(fakeClient, token, version, outputDir)
	})

	Context("when retrieving the file info is successful", func() {
		var (
			privateDownloadURL string
			fileInfo           *slack.File
		)

		BeforeEach(func() {
			privateDownloadURL = fmt.Sprintf("%s/banana.zip", server.URL())
			fileInfo = &slack.File{
				Name:               "banana.zip",
				URLPrivateDownload: privateDownloadURL,
			}
			fakeClient.GetFileInfoReturns(fileInfo, nil, nil, nil)
		})

		Context("when retrieving the file content is successful", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/banana.zip"),
						VerifyHeaderKV("Authorization", "Bearer don't token me bro"),
						RespondWith(http.StatusOK, "I AM A BANANA"),
					),
				)

			})

			It("looks up the file information based on the version provided", func() {
				Expect(executeErr).ToNot(HaveOccurred())

				Expect(fakeClient.GetFileInfoCallCount()).To(Equal(1))
				passedID, _, _ := fakeClient.GetFileInfoArgsForCall(0)
				Expect(passedID).To(Equal(version.ID))
			})

			It("writes the file to disk", func() {

			})
		})

		Context("when retrieving the file content is errors", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/banana.zip"),
						RespondWith(http.StatusNotFound, ""),
					),
				)
			})

			It("returns an error", func() {
				Expect(executeErr).To(MatchError(HTTPError{
					Message: "404 Not Found",
					URL:     privateDownloadURL,
				}))
			})
		})
	})

	Context("when retrieving the file info errors", func() {
		BeforeEach(func() {
			fakeClient.GetFileInfoReturns(nil, nil, nil, errors.New("no, just, no, I can't even"))
		})

		It("returns the error", func() {
			Expect(executeErr).To(MatchError("no, just, no, I can't even"))
		})
	})
})
