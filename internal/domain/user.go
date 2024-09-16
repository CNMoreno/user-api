package domain

import "time"

// User struct of user in BD.
type User struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name" binding:"required" csv:"name" validate:"required"`
	Email     string    `bson:"email" binding:"required,email" csv:"email" validate:"required,email"`
	Enabled   bool      `bson:"enabled"`
	Password  string    `bson:"password" binding:"required,min=8,password" csv:"password" validate:"required,min=8"`
	UserName  string    `bson:"userName" binding:"required" csv:"username" validate:"required"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	DeletedAt time.Time `bson:"deletedAt"`
}
