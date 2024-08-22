package api

import (
	"errors"

	"github.com/marioscordia/auth/pkg/auth_v1"
)

func validateCreateReq(req *auth_v1.CreateRequest) error {
	if req.GetEmail() == "" {
		return errors.New("please fill the email")
	}

	if req.GetName() == "" {
		return errors.New("please fill the username")
	}

	if req.GetPassword() == "" {
		return errors.New("please fill the password")

	}

	if req.GetPassword() != req.GetPasswordConfirm() {
		return errors.New("password is not correctly confirmed")
	}

	return nil
}

func validateUpdateReq(req *auth_v1.UpdateRequest) error {
	if req.GetEmail().GetValue() == "" && req.GetName().GetValue() == "" {
		return errors.New("updating values are empty")
	}

	return nil
}
