package requests

type UpdateUserRequest struct {
	FullName string  `json:"full_name"`
	Login    string  `json:"login"`
	Role     string  `json:"role"`
	Position *string `json:"position"`
}
