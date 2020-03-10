package db

import (
	"github.com/asankov/containerizor/pkg/models"
)

type Database interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)

	Close()
}
