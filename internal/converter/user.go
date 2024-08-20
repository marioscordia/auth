package converter

import (
	"github.com/marioscordia/auth/internal/model"
	"github.com/marioscordia/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserFromService is the method that converts User model to GRPC response
func ToUserFromService(u *model.User) *auth_v1.GetResponse {
	role := auth_v1.Role(auth_v1.Role_value[u.Role])

	return &auth_v1.GetResponse{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      role,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

// ToUserCreateFromCreateRequest is the method that converts GRPC Create request to UserCreate model
func ToUserCreateFromCreateRequest(req *auth_v1.CreateRequest) *model.UserCreate {
	return &model.UserCreate{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

// ToUserUpdateFromUpdateRequest is the method that converts GRPC Update request to UserUpdate model
func ToUserUpdateFromUpdateRequest(req *auth_v1.UpdateRequest) *model.UserUpdate {
	u := &model.UserUpdate{}

	u.ID = req.GetId()

	if req.Email != nil {
		u.Email = req.GetEmail().GetValue()
	}

	if req.Name != nil {
		u.Name = req.GetName().GetValue()
	}

	return u
}
