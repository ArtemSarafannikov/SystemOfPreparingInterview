package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimestampFromPtr(time *time.Time) *timestamppb.Timestamp {
	if time == nil {
		return nil
	}
	return timestamppb.New(*time)
}
