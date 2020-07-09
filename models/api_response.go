package models

// ApiResponse is returned from an api request in the event of an error,
// or if no other data is returned.
type ApiResponse struct {
	// Code returned by the api for use in the client
	Code int32 `json:"code,omitempty"`
	// return type returned by the api for use in the client
	Type_   string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}
