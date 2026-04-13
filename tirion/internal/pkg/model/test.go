package model

import (
	"time"

	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestSchema .
var TestSchema = NewSchema[Test]("problems_tests")

// Test тест для задачи
type Test struct {
	ID         int64     `db:"id"`
	ProblemID  uuid.UUID `db:"problem_id"`
	InputData  string    `db:"input_data"`
	OutputData string    `db:"output_data"`
	Hidden     bool      `db:"is_hidden"`
	UpdatedAt  time.Time `db:"updated_at"`
	CreatedAt  time.Time `db:"created_at"`
}

// ToProto .
func (t Test) ToProto() *tirion.Test {
	return &tirion.Test{
		Id:         t.ID,
		ProblemId:  t.ProblemID.String(),
		InputData:  t.InputData,
		OutputData: t.OutputData,
		IsHidden:   t.Hidden,
		UpdatedAt:  timestamppb.New(t.UpdatedAt),
		CreatedAt:  timestamppb.New(t.CreatedAt),
	}
}
