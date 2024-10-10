package auth

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/taufiksty/be-socket/pkg/domain/entity"
	"github.com/taufiksty/be-socket/pkg/domain/repository"
	"github.com/taufiksty/be-socket/pkg/shared/util"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

func (uc *AuthUsecase) Register(username, email, password string) error {
	userByEmail, err := uc.userRepo.GetUserByEmail(email)
	if (err != nil && err != sql.ErrNoRows) || userByEmail != nil {
		return err
	}
	userByUsername, err := uc.userRepo.GetUserByUsername(username)
	if (err != nil && err != sql.ErrNoRows) || userByUsername != nil {
		return err
	}

	hashedPassword := util.HashPassword(password)

	user := &entity.User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	return uc.userRepo.CreateUser(user)
}

func (uc *AuthUsecase) Login(username, password string) (*entity.User, error) {
	user, err := uc.userRepo.GetUserByUsername(username)
	if err != nil || user == nil {
		return nil, err
	}

	if !util.CheckPassword(password, user.Password) {
		return nil, errors.New("password invalid")
	}

	return user, nil
}
