package repository

import "github.com/taufiksty/be-socket/pkg/domain/entity"

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
}
