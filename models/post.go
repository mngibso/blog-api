package models

type Post struct {
	Id        int64  `json:"id,omitempty"`
	Title     string `json:"title"`
	Status    string `json:"status,omitempty"`
	CreatedAt int32  `json:"createdAt,omitempty"`
	Body      string `json:"body,omitempty"`
}
