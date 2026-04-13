package app

import (
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/store"
	"github.com/CodefriendOrg/tirion/internal/usecase"
)

// Implementation .
type Implementation struct {
	storage *store.Storage
	service *usecase.Service

	tirion.UnimplementedTirionServer
}

// NewImplementation .
func NewImplementation(storage *store.Storage, service *usecase.Service) *Implementation {
	return &Implementation{
		storage: storage,
		service: service,
	}
}
