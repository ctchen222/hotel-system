package models

import (
	"context"

	"github.com/ctchen1999/hotel-system/internal/pgtypes"
)

type PgHotelStore interface {
	CreateHotel(context.Context, *pgtypes.Hotel) error
	GetHotels(context.Context) ([]*pgtypes.Hotel, error)
	GetHotelById(ctx context.Context, id string) (*pgtypes.Hotel, error)
	UpdateHotel(ctx context.Context, hotel *pgtypes.UpdateHotelParams, id string) error
	DeleteHotel(ctx context.Context, id string) error
}

type PostgresHotelStore struct {
	pool *PostgresInstance
}

func NewPostgresHotelStore(pool *PostgresInstance) *PostgresHotelStore {
	return &PostgresHotelStore{
		pool: pool,
	}
}

func (s *PostgresHotelStore) CreateHotel(ctx context.Context, hotel *pgtypes.Hotel) error {
	query := `INSERT INTO hotels(name, location, rating) VALUES($1, $2, $3)`

	_, err := s.pool.DB.Exec(ctx, query, hotel.Name, hotel.Location, hotel.Rating)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresHotelStore) GetHotels(ctx context.Context) ([]*pgtypes.Hotel, error) {
	query := `SELECT * FROM hotels`

	rows, err := s.pool.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []*pgtypes.Hotel
	for rows.Next() {
		var hotel pgtypes.Hotel
		if err := rows.Scan(&hotel.Id, &hotel.Name, &hotel.Location, &hotel.Rating); err != nil {
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}

	return hotels, nil
}

func (s *PostgresHotelStore) GetHotelById(ctx context.Context, id string) (*pgtypes.Hotel, error) {
	query := `SELECT * FROM hotels WHERE id = $1`

	var hotel pgtypes.Hotel
	row := s.pool.DB.QueryRow(ctx, query, id)
	if err := row.Scan(&hotel.Id, &hotel.Name, &hotel.Location, &hotel.Rating); err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (s *PostgresHotelStore) UpdateHotel(ctx context.Context, hotel *pgtypes.UpdateHotelParams, id string) error {
	var hotelId string
	row := s.pool.DB.QueryRow(ctx, `SELECT id FROM hotels WHERE id = $1`, id)
	if err := row.Scan(&hotelId); err != nil {
		return err
	}

	query := `UPDATE hotels SET name = $1, location = $2, rating = $3 WHERE id = $4`
	_, err := s.pool.DB.Exec(ctx, query, hotel.Name, hotel.Location, hotel.Rating, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresHotelStore) DeleteHotel(ctx context.Context, id string) error {
	query := `DELETE FROM hotels WHERE id = $1`

	_, err := s.pool.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
