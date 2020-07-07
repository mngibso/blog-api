package models

type Post struct {
	ID        string `json:"id" bson:"_id"`
	Title     string `json:"title"`
	Status    string `json:"status,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	CreatedBy string `json:"createdBy"`
	Body      string `json:"body"`
}
