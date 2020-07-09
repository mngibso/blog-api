package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// NewPostFromPostInput returns a Post given the id and CreatePostInput object
func NewPostFromPostInput(id interface{}, createPost CreatePostInput) *Post {
	return &Post{
		ID:        id.(primitive.ObjectID),
		Username:  createPost.Username,
		Title:     createPost.Title,
		CreatedAt: createPost.CreatedAt,
		Body:      createPost.Body,
	}
}

type Post struct {
	ID        primitive.ObjectID `json:"id,omitempty" binding:"required" bson:"_id"`
	Title     string             `json:"title" binding:"required"`
	CreatedAt int64              `json:"createdAt,omitempty" bson:"createdAt"`
	Username  string             `json:"username" binding:"required"`
	Body      string             `json:"body" binding:"required"`
}

// CreatePostInput represents post data sent from client
type CreatePostInput struct {
	Title string `json:"title" binding:"required"`
	// seconds since epoch
	CreatedAt int64  `json:"createdAt,omitempty" bson:"createdAt"`
	Username  string `json:"username,omitempty"`
	Body      string `json:"body" binding:"required"`
}
