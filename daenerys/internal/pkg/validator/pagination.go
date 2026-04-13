package validator

import (
	"errors"

	desc "github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	validation "github.com/go-ozzo/ozzo-validation"
)

const maxPerPageLimit = 1000

// ValidatePagination .
func ValidatePagination(value any) error {
	if value == nil {
		return nil
	}

	pagination, ok := value.(*desc.Pagination)
	if !ok {
		return errors.New("value is not pagination")
	}

	if pagination == nil {
		return nil
	}

	return validation.ValidateStruct(pagination,
		validation.Field(&pagination.Page, validation.Required, validation.By(validatePageUint64)),
		validation.Field(&pagination.PerPage, validation.Required, validation.By(validatePerPageUint64)),
	)
}

func validatePageUint64(value any) error {
	page, ok := value.(uint64)
	if !ok {
		return errors.New("value is not uint64")
	}

	pageInt := int64(page) //nolint:gosec

	return validation.Validate(pageInt, validation.Min(1))
}

func validatePerPageUint64(value any) error {
	perPage, ok := value.(uint64)
	if !ok {
		return errors.New("value is not uint64")
	}

	perPageInt := int64(perPage) // nolint:gosec

	return validation.Validate(perPageInt, validation.Max(maxPerPageLimit))
}
