package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/marioscordia/auth/internal/converter"
	"github.com/marioscordia/auth/internal/service"
	"github.com/marioscordia/auth/pkg/auth_v1"
)

// New is a function that returns Handler object
func New(useCase service.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// Handler is an object, which have methods that receive GRPC requests
type Handler struct {
	auth_v1.UnimplementedAuthV1Server
	useCase service.Service
}

// Get is the method that receives GRPC Get request
func (h *Handler) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	user, err := h.useCase.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromService(user), nil
}

// Create is the method that receives GRPC Create request
func (h *Handler) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	if err := validateCreateReq(req); err != nil {
		return nil, err
	}

	u := converter.ToUserCreateFromCreateRequest(req)

	id, err := h.useCase.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return &auth_v1.CreateResponse{Id: id}, nil
}

// Update is the method that receives GRPC Update request
func (h *Handler) Update(ctx context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	if err := validateUpdateReq(req); err != nil {
		return nil, err
	}

	u := converter.ToUserUpdateFromUpdateRequest(req)

	if err := h.useCase.UpdateUser(ctx, u); err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete is the method that receives GRPC Delete request
func (h *Handler) Delete(ctx context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.useCase.DeleteUser(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return nil, nil
}
