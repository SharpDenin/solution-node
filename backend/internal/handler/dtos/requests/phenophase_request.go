package requests

type CreatePhenophaseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	OrderIndex  int    `json:"order_index"`
}

type UpdatePhenophaseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	OrderIndex  int    `json:"order_index"`
}
