package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/asankov/containerizor/pkg/models"

	// to register PostreSQL driver
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func (db *Database) CreateUser(user *models.User) error {
	// TODO: do this in one transaction
	// BIGGER TODO: fix sql injection
	sql := fmt.Sprintf("INSERT INTO USERS(username, passwordHash) VALUES ('%s', '%s');", user.Username, user.HashedPassword)
	if _, err := db.db.Exec(sql); err != nil {
		return err
	}
	sql = fmt.Sprintf("CREATE SCHEMA %s;", user.Username)
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
	connString := fmt.Sprintf("host=%s port=%d user=%s", host, port, user)
	if dbPass != "" {
		connString += fmt.Sprintf(" password=%s", dbPass)
	}
	// apparantly, this must be the last arg to be passed
	connString += fmt.Sprintf(" dbname=%s sslmode=disable", dbName)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := connect(timeoutCtx, connString)
	if err != nil {
		return nil, fmt.Errorf("connect error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	return &Database{
		db: db,
	}, err
}

func connect(ctx context.Context, connString string) (*sql.DB, error) {
	for {
		db, err := sql.Open("postgres", connString)
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err == nil {
			return db, nil
		}

		log.Printf("err: %v", err)
		select {
		case <-time.After(1 * time.Second):
			log.Println("retrying to connect to db")
		case <-ctx.Done():
			log.Println("db connect timeout")
			return nil, errors.New("db connect timeout")
		}
	}
}
