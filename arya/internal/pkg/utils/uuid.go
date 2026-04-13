package utils

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// StringsToUUIDs .
func StringsToUUIDs(strings []string) []uuid.UUID {
	if len(strings) == 0 {
		return nil
	}

	return lo.Map(strings, func(item string, _ int) uuid.UUID {
		return uuid.MustParse(item)
	})
}

// UUIDsToStrings .
func UUIDsToStrings(uuids []uuid.UUID) []string {
	if len(uuids) == 0 {
		return nil
	}

	return lo.Map(uuids, func(item uuid.UUID, _ int) string {
		return item.String()
	})
}
