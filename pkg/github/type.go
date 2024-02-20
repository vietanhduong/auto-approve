package github

type PullRequest struct {
	Number    int
	Title     string
	Body      string
	Author    string
	Reviewers []string
	State     string
	Draft     bool
	Files     []*File
}

type File struct {
	Name   string
	Status string
}

type SubmitEvent string

const (
	SubmitEventApprove        SubmitEvent = "APPROVE"
	SubmitEventRequestChanges SubmitEvent = "REQUEST_CHANGES"
	SubmitEventComment        SubmitEvent = "COMMENT"
)

type SubmitReviewRequest struct {
	Owner  string      // required
	Repo   string      // required
	Number int         // required
	Body   string      // +optional
	Event  SubmitEvent // required
}

func (e SubmitEvent) IsValid() bool {
	switch e {
	case SubmitEventApprove, SubmitEventRequestChanges, SubmitEventComment:
		return true
	}
	return false
}
