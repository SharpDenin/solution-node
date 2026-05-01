package responses

import "time"

type PhenophaseResponse struct {
	ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	Description            *string   `json:"description"`
	ImageURL               *string   `json:"image_url"`
	OrderIndex             int       `json:"order_index"`
	MinCriticalTemperature *float64  `json:"min_critical_temperature"`
	CriticalTemperature    *float64  `json:"critical_temperature"`
	CreatedAt              time.Time `json:"created_at"`
}
