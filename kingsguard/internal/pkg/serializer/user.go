package serializer

import (
	"github.com/CodefriendOrg/kingsguard/internal/model"
	kingsguard "github.com/CodefriendOrg/kingsguard/internal/pb/api"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildUser(user *model.User) *kingsguard.User {
	if user == nil {
		return nil
	}

	return &kingsguard.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		DeletedAt: utils.TimestampFromPtr(user.DeletedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
