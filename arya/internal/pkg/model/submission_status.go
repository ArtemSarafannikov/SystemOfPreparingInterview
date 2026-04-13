package model

import (
	"strings"

	"github.com/CodefriendOrg/arya/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
)

const submissionStatusPrefix = "STATUS_"

// ToProto .
func (e SubmissionStatus) ToProto() daenerys.SubmissionStatus {
	return daenerys.SubmissionStatus(daenerys.SubmissionStatus_value[submissionStatusPrefix+e.String()]) //nolint:unconvert
}

// ConvertSubmissionStatus .
func ConvertSubmissionStatus(status daenerys.SubmissionStatus) SubmissionStatus {
	return SubmissionStatus(strings.TrimPrefix(string(status.String()), submissionStatusPrefix)) //nolint:unconvert
}
