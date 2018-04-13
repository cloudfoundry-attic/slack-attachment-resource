package actor_test

import (
	"errors"

	"github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "code.cloudfoundry.org/slack-attachment-resource/actor"
	fake "code.cloudfoundry.org/slack-attachment-resource/actor/actorfakes"
	"code.cloudfoundry.org/slack-attachment-resource/shared"
)

var _ = Describe("Check", func() {
	var (
		fakeClient *fake.FakeCheckAPIClient
		history    *slack.History

		groupID   string
		filename  string
		timestamp string
		fileID    string

		versions   []shared.Version
		executeErr error
	)

	BeforeEach(func() {
		fakeClient = new(fake.FakeCheckAPIClient)
		history = new(slack.History)

		filename = "our-file"
		groupID = "some-group-id"
		timestamp = ""
		fileID = ""
	})

	JustBeforeEach(func() {
		versions, executeErr = Check(fakeClient, groupID, filename, timestamp)
	})

	Context("when there are no files", func() {
		BeforeEach(func() {
			fakeClient.GetFilesReturns(nil, nil, nil)

			fakeClient.GetGroupHistoryReturns(history, nil)
		})

		It("returns nil", func() {
			Expect(versions).To(BeNil())
			Expect(executeErr).ToNot(HaveOccurred())
		})
	})

	Context("when there are files", func() {
		var (
			files  []slack.File
			getErr error
		)

		Context("when timestamp is empty", func() { //AKA First time check is run
			BeforeEach(func() {
				timestamp = ""

				fakeClient.GetFilesReturns(nil, nil, nil)
			})

			It("calls GetFiles only for the first page", func() {
				Expect(fakeClient.GetFilesCallCount()).To(Equal(1))

				params := fakeClient.GetFilesArgsForCall(0)
				Expect(params).To(MatchFields(IgnoreExtras, Fields{
					"Channel":       Equal(groupID),
					"Page":          BeNumerically("==", 1),
					"TimestampFrom": BeNumerically("==", 0),
					"TimestampTo":   BeNumerically("==", -1),
				}))
			})
		})

		Context("when timestamp is non-empty", func() {
			BeforeEach(func() {
				timestamp = "42"

				files = make([]slack.File, 100)
				fakeClient.GetFilesReturnsOnCall(0, files, nil, nil)

				files = make([]slack.File, 50)
				fakeClient.GetFilesReturnsOnCall(1, files, nil, nil)
			})

			It("calls GetGroupHistory multiple times until all the files since timestamp have been retreived", func() {
				Expect(fakeClient.GetFilesCallCount()).To(Equal(2))

				params1 := fakeClient.GetFilesArgsForCall(0)
				Expect(params1).To(MatchFields(IgnoreExtras, Fields{
					"Channel":       Equal(groupID),
					"Page":          BeNumerically("==", 1),
					"TimestampFrom": BeNumerically("==", 42),
					"TimestampTo":   BeNumerically("==", -1),
				}))

				params2 := fakeClient.GetFilesArgsForCall(1)
				Expect(params2).To(MatchFields(IgnoreExtras, Fields{
					"Channel":       Equal(groupID),
					"Page":          BeNumerically("==", 2),
					"TimestampFrom": BeNumerically("==", 42),
					"TimestampTo":   BeNumerically("==", -1),
				}))
			})
		})

		Context("when multiple files match", func() {
			BeforeEach(func() {
				ourFile := slack.File{Name: "our-file"}
				notFile := slack.File{Name: "not-file"}

				history.Messages = []slack.Message{
					{Msg: slack.Msg{Timestamp: "timestamp-4", File: &ourFile}},
					{Msg: slack.Msg{Timestamp: "timestamp-3"}},
					{Msg: slack.Msg{Timestamp: "timestamp-2", File: &notFile}},
					{Msg: slack.Msg{Timestamp: "timestamp-1", File: &ourFile}},
				}
				getErr = nil
				fakeClient.GetGroupHistoryReturns(history, getErr)

				files = []slack.File{
					slack.File{Name: "our-file", ID: "some-id-3", Created: 3456},
					slack.File{Name: "not-file", ID: "some-id-2", Created: 2345},
					slack.File{Name: "our-file", ID: "some-id-1", Created: 1234},
				}
				fakeClient.GetFilesReturns(files, nil, nil)
			})

			It("returns a list of string versions", func() {
				Expect(versions).To(Equal([]shared.Version{
					{ID: "some-id-1", Timestamp: "1234"},
					{ID: "some-id-3", Timestamp: "3456"},
				}))
				Expect(executeErr).ToNot(HaveOccurred())
			})
		})

		Context("when no files match", func() {
			BeforeEach(func() {
				notFile := slack.File{Name: "not-file"}

				history.Messages = []slack.Message{
					{Msg: slack.Msg{Timestamp: "timestamp-1", File: &notFile}},
					{Msg: slack.Msg{Timestamp: "timestamp-2", File: &notFile}},
					{Msg: slack.Msg{Timestamp: "timestamp-3", File: &notFile}},
				}
				getErr = nil
				fakeClient.GetGroupHistoryReturns(history, getErr)

				files = []slack.File{
					slack.File{Name: "not-file", ID: "some-id-2", Created: 2345},
				}
				fakeClient.GetFilesReturns(files, nil, nil)
			})

			It("returns an empty list of strings", func() {
				Expect(versions).To(BeEmpty())
				Expect(executeErr).ToNot(HaveOccurred())
			})
		})
	})

	Context("when the client returns an errorerrors", func() {
		BeforeEach(func() {
			err := errors.New("some-error")
			fakeClient.GetFilesReturns(nil, nil, err)
		})

		It("raises an error", func() {
			Expect(executeErr).To(MatchError("some-error"))
		})
	})
})
