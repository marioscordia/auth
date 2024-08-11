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
		usecase: useCase,
	}
}

// Handler is an object, which have methods that receive GRPC requests
type Handler struct {
	auth_v1.UnimplementedAuthV1Server
	usecase service.Service
}

// Get is the method that receives GRPC Get request
func (h *Handler) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	user, err := h.usecase.GetUser(ctx, req.Id)
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

	u := converter.ToUserFromCreateRequest(req)

	id, err := h.usecase.CreateUser(ctx, u, req.Password)
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

	u := converter.ToUserFromUpdateRequest(req)

	if err := h.usecase.UpdateUser(ctx, u); err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete is the method that receives GRPC Delete request
func (h *Handler) Delete(ctx context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.usecase.DeleteUser(ctx, req.Id); err != nil {
		return nil, err
	}

	return nil, nil
}
