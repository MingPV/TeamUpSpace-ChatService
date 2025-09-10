package entities

type Order struct {
	ID    uint    `json:"id" bson:"_id,omitempty"`
	Total float64 `json:"total" bson:"total"`
}
