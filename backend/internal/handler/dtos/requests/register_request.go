package requests

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
