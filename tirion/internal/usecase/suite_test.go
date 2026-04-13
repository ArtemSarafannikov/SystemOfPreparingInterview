package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/kingsguard/pkg/kingsguard"
	"github.com/CodefriendOrg/tirion/internal/pkg/store"
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
	var (
		kingsguardClient kingsguard.KingsguardClient
	)

	for _, val := range mocks {
		switch mock := val.(type) {
		case kingsguard.KingsguardClient:
			kingsguardClient = mock
		default:
			s.Failf("Unknown mock type", "unexpected mock type: %T", mock)
		}
	}

	_ = kingsguardClient

	return NewService(s.storage, kingsguardClient)
}

func TestStore(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(testSuite))
}
