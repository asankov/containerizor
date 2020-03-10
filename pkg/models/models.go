package models

import "time"

type User struct {
	ID             int
	Username       string
	Email          string
	HashedPassword string
	Created        time.Time
}
