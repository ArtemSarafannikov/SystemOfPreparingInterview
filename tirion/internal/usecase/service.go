package usecase

import (
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/kingsguard/pkg/kingsguard"
	"github.com/CodefriendOrg/tirion/internal/pkg/store"
)

// Service .
type Service struct {
	storage *store.Storage

	kingsguardClient kingsguard.KingsguardClient
}

// NewService .
func NewService(
	storage *store.Storage,
	kingsguardClient kingsguard.KingsguardClient,
) *Service {
	return &Service{
		storage:          storage,
		kingsguardClient: kingsguardClient,
	}
}
