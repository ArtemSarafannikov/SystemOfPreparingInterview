package app

import (
	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/store"
	"github.com/CodefriendOrg/daenerys/internal/usecase"
)

// Implementation .
type Implementation struct {
	storage *store.Storage
	service *usecase.Service

	daenerys.UnimplementedDaenerysServer
}

// NewImplementation .
func NewImplementation(storage *store.Storage, service *usecase.Service) *Implementation {
	return &Implementation{
		storage: storage,
		service: service,
	}
}
