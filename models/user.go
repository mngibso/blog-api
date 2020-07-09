package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// NewUserFromUserInput returns a User given the id and CreateUserInput object
func NewUserFromUserInput(id interface{}, createUser CreateUserInput) *User {
	return &User{
		ID:        id.(primitive.ObjectID),
		Username:  createUser.Username,
		FirstName: createUser.FirstName,
		LastName:  createUser.LastName,
		Email:     createUser.Email,
	}
}

// User stores information about a user
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" binding:"required" bson:"_id"`
	Username  string             `json:"username" binding:"required"`
	FirstName string             `json:"firstName,omitempty"`
	LastName  string             `json:"lastName,omitempty"`
	Email     string             `json:"email,omitempty"`
	Password  string             `json:"password,omitempty" binding:"required"`
}

// User stores information about a user, sans password
type UserOutput struct {
	ID        primitive.ObjectID `json:"id,omitempty" binding:"required" bson:"_id"`
	Username  string             `json:"username" binding:"required"`
	FirstName string             `json:"firstName,omitempty"`
	LastName  string             `json:"lastName,omitempty"`
	Email     string             `json:"email,omitempty"`
}

// CreateUserInput stores user information obtained from a client request
type CreateUserInput struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty" binding:"required"`
}
