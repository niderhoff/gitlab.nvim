package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harrisoncramer/gitlab.nvim/cmd/app/git"
	"github.com/xanzy/go-gitlab"
)

var errorFromGitlab = errors.New("Some error from Gitlab")

/* The assert function is a helper function used to check two comparables */
func assert[T comparable](t *testing.T, got T, want T) {
	t.Helper()
	if got != want {
		t.Errorf("Got '%v' but wanted '%v'", got, want)
	}
}

/* Will create a new request with the given method, endpoint and body */
func makeRequest(t *testing.T, method string, endpoint string, body any) *http.Request {
	t.Helper()

	var reader io.Reader
	if body != nil {
		j, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(j)
	}

	request, err := http.NewRequest(method, endpoint, reader)
	if err != nil {
		t.Fatal(err)
	}

	return request
}

/* Make response makes a simple response value with the right status code */
func makeResponse(status int) *gitlab.Response {
	return &gitlab.Response{
		Response: &http.Response{
			StatusCode: status,
			Body:       http.NoBody,
		},
	}
}

var testProjectData = data{
	projectInfo: &ProjectInfo{},
	gitInfo: &git.GitData{
		BranchName: "some-branch",
	},
}

func getSuccessData(t *testing.T, svc ServiceWithHandler, request *http.Request) SuccessResponse {
	res := httptest.NewRecorder()
	svc.handler(res, request)

	var data SuccessResponse
	err := json.Unmarshal(res.Body.Bytes(), &data)
	if err != nil {
		t.Error(err)
	}
	return data
}

func getFailData(t *testing.T, svc ServiceWithHandler, request *http.Request) ErrorResponse {
	res := httptest.NewRecorder()
	svc.handler(res, request)

	var data ErrorResponse
	err := json.Unmarshal(res.Body.Bytes(), &data)
	if err != nil {
		t.Error(err)
	}
	return data
}

type testBase struct {
	errFromGitlab bool
	status        int
}

// Helper for easily mocking bad responses or errors from Gitlab
func (f *testBase) handleGitlabError() (*gitlab.Response, error) {
	if f.errFromGitlab {
		return nil, errorFromGitlab
	}
	if f.status == 0 {
		f.status = 200
	}
	return makeResponse(f.status), nil
}

func checkErrorFromGitlab(t *testing.T, data ErrorResponse, msg string) {
	t.Helper()
	assert(t, data.Status, http.StatusInternalServerError)
	assert(t, data.Message, msg)
	assert(t, data.Details, errorFromGitlab.Error())
}

func checkBadMethod(t *testing.T, data ErrorResponse, methods ...string) {
	t.Helper()
	assert(t, data.Status, http.StatusMethodNotAllowed)
	assert(t, data.Details, "Invalid request type")
	expectedMethods := strings.Join(methods, " or ")
	assert(t, data.Message, fmt.Sprintf("Expected %s", expectedMethods))
}

func checkNon200(t *testing.T, data ErrorResponse, msg, endpoint string) {
	t.Helper()
	assert(t, data.Status, http.StatusSeeOther)
	assert(t, data.Message, msg)
	assert(t, data.Details, fmt.Sprintf("An error occurred on the %s endpoint", endpoint))
}

type FakeGitManager struct {
	RemoteUrl   string
	BranchName  string
	ProjectName string
	Namespace   string
}

func (f FakeGitManager) RefreshProjectInfo(remote string) error {
	return nil
}

func (f FakeGitManager) GetCurrentBranchNameFromNativeGitCmd() (string, error) {
	return f.BranchName, nil
}

func (f FakeGitManager) GetLatestCommitOnRemote(remote string, branchName string) (string, error) {
	return "", nil
}

func (f FakeGitManager) GetProjectUrlFromNativeGitCmd(string) (url string, err error) {
	return f.RemoteUrl, nil
}
