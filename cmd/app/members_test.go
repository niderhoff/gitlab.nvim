package app

import (
	"net/http"
	"testing"

	"github.com/xanzy/go-gitlab"
)

type fakeMemberLister struct {
	testBase
}

func (f fakeMemberLister) ListAllProjectMembers(pid interface{}, opt *gitlab.ListProjectMembersOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.ProjectMember, *gitlab.Response, error) {
	resp, err := f.handleGitlabError()
	if err != nil {
		return nil, nil, err
	}
	return []*gitlab.ProjectMember{}, resp, err
}

func TestMembersHandler(t *testing.T) {
	t.Run("Returns project members", func(t *testing.T) {
		request := makeRequest(t, http.MethodGet, "/project/members", nil)
		svc := projectMemberService{testProjectData, fakeMemberLister{}}
		data := getSuccessData(t, svc, request)
		assert(t, data.Status, http.StatusOK)
		assert(t, data.Message, "Project members retrieved")
	})
	t.Run("Disallows non-GET methods", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/project/members", nil)
		svc := projectMemberService{testProjectData, fakeMemberLister{}}
		data := getFailData(t, svc, request)
		checkBadMethod(t, data, http.MethodGet)
	})
	t.Run("Handles error from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodGet, "/project/members", nil)
		svc := projectMemberService{testProjectData, fakeMemberLister{testBase{errFromGitlab: true}}}
		data := getFailData(t, svc, request)
		checkErrorFromGitlab(t, data, "Could not retrieve project members")
	})
	t.Run("Handles non-200s from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodGet, "/project/members", nil)
		svc := projectMemberService{testProjectData, fakeMemberLister{testBase{status: http.StatusSeeOther}}}
		data := getFailData(t, svc, request)
		checkNon200(t, data, "Could not retrieve project members", "/project/members")
	})
}
