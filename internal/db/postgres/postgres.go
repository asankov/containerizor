package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/asankov/containerizor/pkg/models"

	// to register PostreSQL driver
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func (db *Database) CreateUser(user *models.User) error {
	// TODO: do this in one transaction
	if _, err := db.db.Exec(`INSERT INTO USERS(username, passwordHash) VALUES ($1, $2);`, user.Username, user.HashedPassword); err != nil {
		return err
	}
	// the PostgreSQL driver refuses to format this, so we must go with fmt.Sprintf
	sql := fmt.Sprintf("CREATE SCHEMA %s;", user.Username)
	if _, err := db.db.Exec(sql); err != nil {
		return err
	}
	// TODO: this could be template
	sql = fmt.Sprintf(`CREATE TABLE %s.CONTAINERS (id TEXT)`, user.Username)
	if _, err := db.db.Exec(sql); err != nil {
		return err
	}
	return nil
}
func (db *Database) GetUserByID(id int) (*models.User, error) {
	return nil, nil
}
func (db *Database) GetUserByUsername(username string) (*models.User, error) {
	user := new(models.User)

	err := db.db.
		QueryRow("SELECT * FROM USERS U WHERE U.USERNAME = $1", username).
		Scan(&user.ID, &user.Username, &user.HashedPassword)
	if err != nil {
		// TODO: maybe return proper error here
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (db *Database) Close() {
	db.Close()
}

// New connects to a PostgreSQL instance with the
// given parameters and returns the connection,
// or an error if such occured.
func New(host string, port int, user string, dbName string, dbPass string) (*Database, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s", host, port, user, dbName)
	if dbPass != "" {
		connString += fmt.Sprintf(" password=%s", dbPass)
	}
	// apparantly, this must be the last arg to be passed
	connString += " sslmode=disable"
	db, err := sql.Open("postgres", connString)
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
