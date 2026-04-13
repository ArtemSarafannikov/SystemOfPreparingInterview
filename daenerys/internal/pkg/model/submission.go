package model

import (
	"time"

	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SubmissionSchema .
var SubmissionSchema = NewSchema[Submission]("submissions")

// Submission .
type Submission struct {
	ID        uuid.UUID                    `db:"id"`
	ProblemID uuid.UUID                    `db:"problem_id"`
	UserID    uuid.UUID                    `db:"user_id"`
	Status    daenerys.SubmissionStatus    `db:"status"`
	Code      string                       `db:"code"`
	Language  daenerys.ProgrammingLanguage `db:"language"`
	UpdatedAt time.Time                    `db:"updated_at"`
	CreatedAt time.Time                    `db:"created_at"`
}

// ToProto .
func (s Submission) ToProto() *daenerys.Submission {
	return &daenerys.Submission{
		Id:        s.ID.String(),
		ProblemId: s.ProblemID.String(),
		UserId:    s.UserID.String(),
		Status:    s.Status,
		Code:      s.Code,
		Language:  s.Language,
		UpdatedAt: timestamppb.New(s.UpdatedAt),
		CreatedAt: timestamppb.New(s.CreatedAt),
	}
}
