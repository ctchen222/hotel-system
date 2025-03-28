package pgtypes

type Hotel struct {
	Id       int    `db:"id,omitempty" json:"id,omitempty"`
	Name     string `db:"name" json:"name,omitempty"`
	Location string `db:"location" json:"location,omitempty"`
	Rating   int    `db:"rating" json:"rating,omitempty"`
}

type CreateHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}

type UpdateHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}
