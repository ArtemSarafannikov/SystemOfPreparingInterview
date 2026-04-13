package graph

import (
	"math"

	"github.com/CodefriendOrg/arya/internal/pkg/helper/daenerys_helper"
	"github.com/CodefriendOrg/arya/internal/pkg/helper/kingsguard_helper"
	"github.com/CodefriendOrg/arya/internal/pkg/helper/tirion_helper"
	"github.com/CodefriendOrg/arya/internal/pkg/model"
)

// Resolver .
type Resolver struct {
	kingsguardHelper *kingsguard_helper.Service
	daenerysHelper   *daenerys_helper.Service
	tirionHelper     *tirion_helper.Service
}

// NewResolver .
func NewResolver(
	kingsguardHelper *kingsguard_helper.Service,
	daenerysHelper *daenerys_helper.Service,
	tirionHelper *tirion_helper.Service,
) *Resolver {
	return &Resolver{
		kingsguardHelper: kingsguardHelper,
		daenerysHelper:   daenerysHelper,
		tirionHelper:     tirionHelper,
	}
}

func (r *Resolver) buildPaginationPageInfo(page int64, perPage int64, totalItems uint64) *model.PaginationInfo {
	totalPages := int64(math.Ceil(float64(totalItems) / float64(perPage)))
	if totalPages == 0 && totalItems != 0 {
		totalPages = 1
	}
	hasNextPage := false
	if page < totalPages {
		hasNextPage = true
	}
	hasPreviousPage := false
	if page > 1 && page <= (totalPages+1) { // TODO: спорный момент
		hasPreviousPage = true
	}

	return &model.PaginationInfo{
		TotalPages:      totalPages,
		TotalItems:      int64(totalItems),
		Page:            page,
		PerPage:         perPage,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}
}
