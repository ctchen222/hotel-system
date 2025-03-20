package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	RoomId    primitive.ObjectID `bson:"roomId,omitempty" json:"roomId,omitempty"`
	NumPerson int                `bson:"numPerson,omitempty" json:"numPerson,omitempty"`
	From      time.Time          `bson:"from,omitempty" json:"from,omitempty"`
	To        time.Time          `bson:"to,omitempty" json:"to,omitempty"`
}

type BookingParams struct {
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	NumPerson int       `json:"numPerson"`
}

func (p *BookingParams) Validate() map[string]string {
	now := time.Now()
	errors := map[string]string{}
	if now.After(p.From) {
		errors["from"] = fmt.Sprintf("Can't book room in the past")
	}
	if now.After(p.To) {
		errors["to"] = fmt.Sprintf("Can't book room in the past")
	}
	if p.From.After(p.To) {
		errors["order"] = fmt.Sprintf("From Date After To Date")
	}
	return errors
}

type BookingRawParams struct {
	From      string `json:"from"`
	To        string `json:"to"`
	NumPerson int    `json:"numPerson"`
}

type BookingQuery struct {
	From time.Time
	To   time.Time
}
