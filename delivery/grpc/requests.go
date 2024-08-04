package grpc

import (
	"errors"

	"github.com/marioscordia/auth/pkg/auth_v1"
)

func validateCreateReq(req *auth_v1.CreateRequest) error {
	if req.Email == "" {
		return errors.New("please fill the email")
	}

	if req.Name == "" {
		return errors.New("please fill the username")
	}

	if req.Password == "" {
		return errors.New("please fill the password")

	}

	if req.Password != req.PasswordConfirm {
		return errors.New("password is not correctly confirmed")
	}

	return nil
}

func validateUpdateReq(req *auth_v1.UpdateRequest) error {
	if req.Email.Value == "" && req.Name.Value == "" {
		return errors.New("updating values are empty")
	}

	return nil
}
