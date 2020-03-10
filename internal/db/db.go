package db

import (
	"github.com/asankov/containerizor/pkg/models"
)

type Database interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUsernameAndPassword(username, password string) (*models.User, error)
	Close()
}
