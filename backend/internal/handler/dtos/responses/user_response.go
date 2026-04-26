package responses

type CurrentUserResponse struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	Login    string  `json:"login"`
	Role     string  `json:"role"`
	Position *string `json:"position"`
}
