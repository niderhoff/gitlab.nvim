package app

import (
	"net/http"
	"testing"

	"github.com/xanzy/go-gitlab"
)

type fakeCommentClient struct {
	testBase
}

func (f fakeCommentClient) CreateMergeRequestDiscussion(pid interface{}, mergeRequest int, opt *gitlab.CreateMergeRequestDiscussionOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Discussion, *gitlab.Response, error) {
	resp, err := f.handleGitlabError()
	if err != nil {
		return nil, nil, err
	}

	return &gitlab.Discussion{Notes: []*gitlab.Note{{}}}, resp, err
}
func (f fakeCommentClient) UpdateMergeRequestDiscussionNote(pid interface{}, mergeRequest int, discussion string, note int, opt *gitlab.UpdateMergeRequestDiscussionNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error) {
	resp, err := f.handleGitlabError()
	if err != nil {
		return nil, nil, err
	}

	return &gitlab.Note{}, resp, err
}

func (f fakeCommentClient) DeleteMergeRequestDiscussionNote(pid interface{}, mergeRequest int, discussion string, note int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error) {
	resp, err := f.handleGitlabError()
	if err != nil {
		return nil, err
	}
	return resp, err
}

func TestPostComment(t *testing.T) {
	var testCommentCreationData = PostCommentRequest{Comment: "Some comment"}
	t.Run("Creates a new note (unlinked comment)", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/mr/comment", testCommentCreationData)
		svc := commentService{testProjectData, fakeCommentClient{}}
		data := getSuccessData(t, svc, request)
		assert(t, data.Message, "Comment created successfully")
		assert(t, data.Status, http.StatusOK)
	})

	t.Run("Creates a new comment", func(t *testing.T) {
		testCommentCreationData := PostCommentRequest{ // Re-create comment creation data to avoid mutating this variable in other tests
			Comment: "Some comment",
			PositionData: PositionData{
				FileName: "file.txt",
			},
		}
		request := makeRequest(t, http.MethodPost, "/mr/comment", testCommentCreationData)
		svc := commentService{testProjectData, fakeCommentClient{}}
		data := getSuccessData(t, svc, request)
		assert(t, data.Message, "Comment created successfully")
		assert(t, data.Status, http.StatusOK)
	})

	t.Run("Handles errors from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/mr/comment", testCommentCreationData)
		svc := commentService{testProjectData, fakeCommentClient{testBase{errFromGitlab: true}}}
		data := getFailData(t, svc, request)
		checkErrorFromGitlab(t, data, "Could not create discussion")
	})

	t.Run("Handles non-200s from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/mr/comment", testCommentCreationData)
		svc := commentService{testProjectData, fakeCommentClient{testBase{status: http.StatusSeeOther}}}
		data := getFailData(t, svc, request)
		checkNon200(t, data, "Could not create discussion", "/mr/comment")
	})
}

func TestDeleteComment(t *testing.T) {
	var testCommentDeletionData = DeleteCommentRequest{NoteId: 3, DiscussionId: "abc123"}
	t.Run("Deletes a comment", func(t *testing.T) {
		request := makeRequest(t, http.MethodDelete, "/mr/comment", testCommentDeletionData)
		svc := commentService{testProjectData, fakeCommentClient{}}
		data := getSuccessData(t, svc, request)
		assert(t, data.Message, "Comment deleted successfully")
		assert(t, data.Status, http.StatusOK)
	})
	t.Run("Handles errors from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodDelete, "/mr/comment", testCommentDeletionData)
		svc := commentService{testProjectData, fakeCommentClient{testBase{errFromGitlab: true}}}
		data := getFailData(t, svc, request)
		checkErrorFromGitlab(t, data, "Could not delete comment")
	})

	t.Run("Handles non-200s from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodDelete, "/mr/comment", testCommentDeletionData)
		svc := commentService{testProjectData, fakeCommentClient{testBase{status: http.StatusSeeOther}}}
		data := getFailData(t, svc, request)
		checkNon200(t, data, "Could not delete comment", "/mr/comment")
	})
}

func TestEditComment(t *testing.T) {
	var testEditCommentData = EditCommentRequest{Comment: "Some comment", NoteId: 3, DiscussionId: "abc123"}
	t.Run("Edits a comment", func(t *testing.T) {
		request := makeRequest(t, http.MethodPatch, "/mr/comment", testEditCommentData)
		svc := commentService{testProjectData, fakeCommentClient{}}
		data := getSuccessData(t, svc, request)
		assert(t, data.Message, "Comment updated successfully")
		assert(t, data.Status, http.StatusOK)
	})
	t.Run("Handles errors from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPatch, "/mr/comment", testEditCommentData)
		svc := commentService{testProjectData, fakeCommentClient{testBase{errFromGitlab: true}}}
		data := getFailData(t, svc, request)
		checkErrorFromGitlab(t, data, "Could not update comment")
	})

	t.Run("Handles non-200s from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPatch, "/mr/comment", testEditCommentData)
		svc := commentService{testProjectData, fakeCommentClient{testBase{status: http.StatusSeeOther}}}
		data := getFailData(t, svc, request)
		checkNon200(t, data, "Could not update comment", "/mr/comment")
	})
}
