package app

import (
	"net/http"
	"testing"

	"github.com/xanzy/go-gitlab"
)

type fakeReplyManager struct {
	testBase
}

func (f fakeReplyManager) AddMergeRequestDiscussionNote(interface{}, int, string, *gitlab.AddMergeRequestDiscussionNoteOptions, ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error) {
	resp, err := f.handleGitlabError()
	if err != nil {
		return nil, nil, err
	}

	return &gitlab.Note{}, resp, err
}

func TestReplyHandler(t *testing.T) {
	var testReplyRequest = ReplyRequest{DiscussionId: "abc123", Reply: "Some Reply", IsDraft: false}
	t.Run("Sends a reply", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/mr/reply", testReplyRequest)
		svc := replyService{testProjectData, fakeReplyManager{}}
		data := getSuccessData(t, svc, request)
		assert(t, data.Message, "Replied to comment")
		assert(t, data.Status, http.StatusOK)
	})
	t.Run("Handles errors from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/mr/reply", testReplyRequest)
		svc := replyService{testProjectData, fakeReplyManager{testBase{errFromGitlab: true}}}
		data := getFailData(t, svc, request)
		checkErrorFromGitlab(t, data, "Could not leave reply")
	})

	t.Run("Handles non-200s from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/mr/reply", testReplyRequest)
		svc := replyService{testProjectData, fakeReplyManager{testBase{status: http.StatusSeeOther}}}
		data := getFailData(t, svc, request)
		checkNon200(t, data, "Could not leave reply", "/mr/reply")
	})
}
