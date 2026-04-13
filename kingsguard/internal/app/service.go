package app

import (
	kingsguard "github.com/CodefriendOrg/kingsguard/internal/pb/api"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/store"
	"github.com/CodefriendOrg/kingsguard/internal/usecase"
)

type Implementation struct {
	storage *store.Storage
	service *usecase.Service

	kingsguard.UnimplementedKingsguardServer
}

func NewService(storage *store.Storage, service *usecase.Service) *Implementation {
	return &Implementation{
		storage: storage,
		service: service,
	}
}
