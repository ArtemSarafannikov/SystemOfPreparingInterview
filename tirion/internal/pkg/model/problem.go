package model

import (
	"time"

	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProblemSchema .
var ProblemSchema = NewSchema[Problem]("problems")

// Problem .
type Problem struct {
	ID            uuid.UUID `db:"id"`
	Summary       string    `db:"summary"`
	Description   string    `db:"description"`
	AuthorID      uuid.UUID `db:"author_id"`
	TimeLimitMs   int64     `db:"time_limit_ms"`
	MemoryLimitKb int64     `db:"memory_limit_kb"`
	UpdatedAt     time.Time `db:"updated_at"`
	CreatedAt     time.Time `db:"created_at"`
}

// ToProto .
func (t Problem) ToProto() *tirion.Problem {
	return &tirion.Problem{
		Id:            t.ID.String(),
		Summary:       t.Summary,
		Description:   t.Description,
		AuthorId:      t.AuthorID.String(),
		TimeLimitMs:   t.TimeLimitMs,
		MemoryLimitKb: t.MemoryLimitKb,
		UpdatedAt:     timestamppb.New(t.UpdatedAt),
		CreatedAt:     timestamppb.New(t.CreatedAt),
	}
}
