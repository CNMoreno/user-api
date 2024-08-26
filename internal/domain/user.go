package domain

import "time"

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Enabled   bool      `bson:"enabled"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DelatedAt time.Time `bson:"deleted_at"`
}
