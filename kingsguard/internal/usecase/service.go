package usecase

import "github.com/CodefriendOrg/kingsguard/internal/pkg/store"

type Service struct {
	storage *store.Storage
}

func NewService(storage *store.Storage) *Service {
	return &Service{
		storage: storage,
	}
}
