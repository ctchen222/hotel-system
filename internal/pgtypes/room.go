package pgtypes

type Room struct {
	Id      int     `db:"id,omitempty" json:"id,omitempty"`
	Size    string  `db:"size" json:"size"`
	SeaSide bool    `db:"seaside" json:"seaside"`
	Price   float64 `db:"price" json:"price"`
	HotelId int     `db:"hotelId" json:"hotelId"`
}

type CreateRoomParams struct {
	Size    string  `json:"size,omitempty"`
	SeaSide bool    `json:"seaside,omitempty"`
	Price   float64 `json:"price,omitempty"`
}
