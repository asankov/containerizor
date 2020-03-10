package postgres

import (
	"database/sql"
	"fmt"

	"github.com/asankov/containerizor/pkg/models"
)

type Database struct {
	db *sql.DB
}

func (db *Database) CreateUser(user *models.User) (*models.User, error) {
	return nil, nil
}
func (db *Database) GetUserByID(id int) (*models.User, error) {
	return nil, nil
}
func (db *Database) GetUserByUsernameAndPassword(username, password string) (*models.User, error) {
	return nil, nil
}
func (db *Database) Close() {
	db.Close()
}

// New connects to a PostgreSQL instance with the
// given parameters and returns the connection,
// or an error if such occured.
func New(host string, port int, user string, dbName string) (*Database, error) {
	// TODO: password, sslmode
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbName))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, err
}
