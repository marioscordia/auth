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

// ToUserFromCreateRequest is the method that converts GRPC Create request to User model
func ToUserFromCreateRequest(req *auth_v1.CreateRequest) *model.User {
	return &model.User{
		Name:  req.Name,
		Email: req.Email,
	}
}

// ToUserFromUpdateRequest is the method that converts GRPC Update request to User model
func ToUserFromUpdateRequest(req *auth_v1.UpdateRequest) *model.User {
	u := &model.User{}

	u.ID = req.Id

	if req.Email != nil {
		u.Email = req.Name.Value
	}

	if req.Name != nil {
		u.Name = req.Name.Value
	}

	return u
}
