package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/marioscordia/auth/pkg/auth_v1"
	"github.com/marioscordia/auth/service"
)

// New is ...
func New(useCase service.Service) *Handler {
	return &Handler{
		usecase: useCase,
	}
}

// Handler is ...
type Handler struct {
	auth_v1.UnimplementedAuthV1Server
	usecase service.Service
}

// Get is ...
func (h *Handler) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	return h.usecase.GetUser(ctx, req.Id)
}

// Create is ...
func (h *Handler) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	if err := validateCreateReq(req); err != nil {
		return nil, err
	}

	id, err := h.usecase.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &auth_v1.CreateResponse{Id: id}, nil
}

// Update is ...
func (h *Handler) Update(ctx context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	if err := validateUpdateReq(req); err != nil {
		return nil, err
	}

	if err := h.usecase.UpdateUser(ctx, req); err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete is ...
func (h *Handler) Delete(ctx context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.usecase.DeleteUser(ctx, req.Id); err != nil {
		return nil, err
	}

	return nil, nil
}
