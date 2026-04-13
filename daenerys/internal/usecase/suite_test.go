package usecase

import (
	"context"
	"testing"

	"github.com/CodefriendOrg/daenerys/internal/pkg/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite

	ctx     context.Context
	ctrl    *gomock.Controller
	storage *store.Storage
}

func (s *testSuite) SetupSuite() {
	storage, err := store.NewTest(context.Background())
	s.Require().NoError(err)
	s.storage = storage
}

func (s *testSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
}

func (s *testSuite) newMockedService(mocks ...any) *Service {
	for _, val := range mocks {
		switch mock := val.(type) {
		default:
			s.Failf("Unknown mock type", "unexpected mock type: %T", mock)
		}
	}

	return NewService(s.storage, nil)
}

func TestStore(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(testSuite))
}
