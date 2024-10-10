package controller

import (
	"github.com/taufiksty/be-socket/pkg/shared/util"
	"github.com/taufiksty/be-socket/pkg/usecase/auth"
)

type UserController struct {
	authUsecase auth.AuthUsecase
}

func NewUserController(authUsecase *auth.AuthUsecase) *UserController {
	return &UserController{authUsecase: *authUsecase}
}

func (ctrl *UserController) Register(username, email, password string) error {
	err := ctrl.authUsecase.Register(username, email, password)
	if err != nil {
		return err
	}
	return nil
}

func (ctrl *UserController) Login(username, password string) (interface{}, error) {
	user, err := ctrl.authUsecase.Login(username, password)
	if err != nil {
		return nil, err
	}

	token, err := util.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
		"user":  user,
	}, nil
}
