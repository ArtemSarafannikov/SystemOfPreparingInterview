package python

import "github.com/google/uuid"

// Args .
type Args struct {
	SubmissionID  uuid.UUID `json:"submission_id"`
	ProblemID     uuid.UUID `json:"problem_id"`
	Code          string    `json:"code"`
	PythonVersion string    `json:"python_version"`
}

// Kind .
func (Args) Kind() string {
	return "judge-submission-python"
}
