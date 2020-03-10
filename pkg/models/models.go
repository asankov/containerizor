package models

import "time"

type User struct {
	ID             int
	Username       string
	HashedPassword string
	Created        time.Time
}
