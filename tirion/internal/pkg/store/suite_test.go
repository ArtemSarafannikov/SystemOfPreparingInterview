package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite

	ctx     context.Context
	storage *Storage
}

func (s *testSuite) SetupSuite() {
	storage, err := NewTest(context.Background())
	s.Require().NoError(err)
	s.storage = storage
}

func (s *testSuite) SetupTest() {
	s.ctx = context.Background()
}

func TestStore(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(testSuite))
}
