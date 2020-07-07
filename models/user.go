package models

type User struct {
	ID        string `json:"id,omitempty" bson:"_id"`
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty" binding:"required"`
}

type CreateUserInput struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty" binding:"required"`
}
