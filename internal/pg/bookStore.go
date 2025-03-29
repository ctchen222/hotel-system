package models

import (
	"context"

	"github.com/ctchen1999/hotel-system/internal/pgtypes"
)

type BookingStore interface {
	CreateBooking(context.Context, *pgtypes.Booking) error
	GetBookingByUserId(ctx context.Context, userId string) ([]*pgtypes.BookingInfo, error)
}

type PostgresBookingStore struct {
	pool *PostgresInstance
}

func NewPostgresBookingStore(pool *PostgresInstance) *PostgresBookingStore {
	return &PostgresBookingStore{
		pool: pool,
	}
}

func (s *PostgresBookingStore) CreateBooking(ctx context.Context, booking *pgtypes.Booking) error {
	query := `INSERT INTO 
		bookings (userid, roomid, numperson, fromdate, todate) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err := s.pool.DB.Exec(ctx, query,
		booking.UserId,
		booking.RoomId,
		booking.NumPerson,
		booking.FromDate,
		booking.ToDate)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresBookingStore) GetBookingByUserId(ctx context.Context, userId string) ([]*pgtypes.BookingInfo, error) {
	query := `SELECT u.id, u.firstname, u.email, b.numperson
		FROM users u
		LEFT JOIN bookings b
		ON u.id = b.userid
		WHERE u.id = $1`

	rows, err := s.pool.DB.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	var bookingInfos []*pgtypes.BookingInfo
	for rows.Next() {
		var bookingInfo pgtypes.BookingInfo
		if err := rows.Scan(&bookingInfo.Id, &bookingInfo.Firstname, &bookingInfo.Email, &bookingInfo.NumPerson); err != nil {
			return nil, err
		}

		bookingInfos = append(bookingInfos, &bookingInfo)
	}

	return bookingInfos, nil
}
