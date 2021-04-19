package entities

import "time"

type Member struct {
	ID           int
	Username     string
	Avatar       string
	CreatedAt    time.Time `db:"created_at"`
	PasswordHash string    `db:"password_hash"`
	Following    int
	FollowedBy   int `db:"followed_by"`
}
