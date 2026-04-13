package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
)

func Test_ValidatePagination(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		t.Parallel()

		err := ValidatePagination(nil)
		require.NoError(t, err)
	})

	t.Run("value is not pagination", func(t *testing.T) {
		t.Parallel()

		err := ValidatePagination(1)
		require.EqualError(t, err, "value is not pagination")
	})

	t.Run("value is pagination, but nil", func(t *testing.T) {
		t.Parallel()

		err := ValidatePagination((*tirion.Pagination)(nil))
		require.NoError(t, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		err := ValidatePagination(&tirion.Pagination{
			Page:    1,
			PerPage: 10,
		})
		require.NoError(t, err)
	})
}
