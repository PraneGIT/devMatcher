package dto

// UserProfileResponse is used for returning user profile data
// (expand as needed for more fields)
type UserProfileResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserProfileUpdateRequest is used for updating user profile
// (expand as needed for more fields)
type UserProfileUpdateRequest struct {
	Name string `json:"name"`
}
