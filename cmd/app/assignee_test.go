package app

import (
	"net/http"
	"testing"

	"github.com/xanzy/go-gitlab"
)

type fakeAssigneeClient struct {
	testBase
}

func (f fakeAssigneeClient) UpdateMergeRequest(pid interface{}, mergeRequest int, opt *gitlab.UpdateMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	resp, err := f.handleGitlabError()
	if err != nil {
		return nil, nil, err
	}
	return &gitlab.MergeRequest{}, resp, nil
}

func TestAssigneeHandler(t *testing.T) {
	var updatePayload = AssigneeUpdateRequest{Ids: []int{1, 2}}

	t.Run("Updates assignees", func(t *testing.T) {
		request := makeRequest(t, http.MethodPut, "/mr/assignee", updatePayload)
		client := fakeAssigneeClient{}
		svc := assigneesService{testProjectData, client}
		data := getSuccessData(t, svc, request)
		assert(t, data.Message, "Assignees updated")
		assert(t, data.Status, http.StatusOK)
	})

	t.Run("Disallows non-PUT method", func(t *testing.T) {
		request := makeRequest(t, http.MethodGet, "/mr/assignee", nil)
		client := fakeAssigneeClient{}
		svc := assigneesService{testProjectData, client}
		data := getFailData(t, svc, request)
		checkBadMethod(t, data, http.MethodPut)
	})

	t.Run("Handles errors from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPut, "/mr/approve", updatePayload)
		client := fakeAssigneeClient{testBase{errFromGitlab: true}}
		svc := assigneesService{testProjectData, client}
		data := getFailData(t, svc, request)
		checkErrorFromGitlab(t, data, "Could not modify merge request assignees")
	})

	t.Run("Handles non-200s from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPut, "/mr/approve", updatePayload)
		client := fakeAssigneeClient{testBase{status: http.StatusSeeOther}}
		svc := assigneesService{testProjectData, client}
		data := getFailData(t, svc, request)
		checkNon200(t, data, "Could not modify merge request assignees", "/mr/assignee")
	})
}
