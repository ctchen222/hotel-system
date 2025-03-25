package types

import (
	"reflect"
	"testing"
	"time"
)

func TestBookingParams_Validate(t *testing.T) {
	type fields struct {
		From      time.Time
		To        time.Time
		NumPerson int
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "Valid Booking Params",
			fields: fields{
				From:      time.Now().Add(time.Hour * 24),
				To:        time.Now().Add(time.Hour * 48),
				NumPerson: 2,
			},
			want: map[string]string{},
		},
		{
			name: "Invalid Booking Params",
			fields: fields{
				From:      time.Now().Add(-time.Hour * 48),
				To:        time.Now().Add(-time.Hour * 24),
				NumPerson: 2,
			},
			want: map[string]string{
				"from": "Can't book room in the past",
				"to":   "Can't book room in the past",
			},
		},
		{
			name: "Invalid Booking Params",
			fields: fields{
				From:      time.Now().Add(-time.Hour * 24),
				To:        time.Now().Add(-time.Hour * 48),
				NumPerson: 3,
			},
			want: map[string]string{
				"from":  "Can't book room in the past",
				"to":    "Can't book room in the past",
				"order": "From Date After To Date",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &BookingParams{
				From:      tt.fields.From,
				To:        tt.fields.To,
				NumPerson: tt.fields.NumPerson,
			}
			if got := p.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BookingParams.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
