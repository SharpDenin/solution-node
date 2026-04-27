package responses

type CurrentUserResponse struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	Login    string  `json:"login"`
	Role     string  `json:"role"`
	Position *string `json:"position"`
}

type UserAdminResponse struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	Login     string  `json:"login"`
	Role      string  `json:"role"`
	Position  *string `json:"position"`
	IsBlocked bool    `json:"is_blocked"`
	IsDeleted bool    `json:"is_deleted"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
