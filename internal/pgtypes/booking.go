package pgtypes

import "time"

type Booking struct {
	Id        int       `db:"id" json:"id,omitempty"`
	UserId    int       `db:"userid" json:"userId,omitempty"`
	RoomId    int       `db:"roomid" json:"roomId,omitempty"`
	NumPerson int       `db:"numperson" json:"numperson,omitempty"`
	FromDate  time.Time `db:"fromdate" json:"fromdate,omitempty"`
	ToDate    time.Time `db:"todate" json:"todate,omitempty"`
}

type BookingParams struct {
	RoomId    string `json:"roomId"`
	FromDate  string `json:"fromdate"`
	ToDate    string `json:"todate"`
	NumPerson int    `json:"numperson"`
}

type BookingInfo struct {
	Id        int    `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Email     string `json:"email,omitempty"`
	NumPerson int    `json:"numperson,omitempty"`
}
