package requests

type CreateVarietyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	ImageURL    string `json:"image_url"`
}

type UpdateVarietyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	ImageURL    string `json:"image_url"`
}
