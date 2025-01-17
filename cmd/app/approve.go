package app

import (
	"encoding/json"
	"net/http"

	"github.com/xanzy/go-gitlab"
)

type MergeRequestApprover interface {
	ApproveMergeRequest(pid interface{}, mr int, opt *gitlab.ApproveMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequestApprovals, *gitlab.Response, error)
}

type mergeRequestApproverService struct {
	data
	client MergeRequestApprover
}

/* approveHandler approves a merge request. */
func (a mergeRequestApproverService) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.Header().Set("Access-Control-Allow-Methods", http.MethodPost)
		handleError(w, InvalidRequestError{}, "Expected POST", http.StatusMethodNotAllowed)
		return
	}

	_, res, err := a.client.ApproveMergeRequest(a.projectInfo.ProjectId, a.projectInfo.MergeId, nil, nil)

	if err != nil {
		handleError(w, err, "Could not approve merge request", http.StatusInternalServerError)
		return
	}

	if res.StatusCode >= 300 {
		handleError(w, GenericError{endpoint: "/mr/approve"}, "Could not approve merge request", res.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := SuccessResponse{
		Message: "Approved MR",
		Status:  http.StatusOK,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		handleError(w, err, "Could not encode response", http.StatusInternalServerError)
	}
}
