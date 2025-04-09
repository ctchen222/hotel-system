package models

import (
	"context"

	"github.com/ctchen222/hotel-system/internal/pgtypes"
)

type PgRoomStore interface {
	CreateRoom(ctx context.Context, room pgtypes.CreateRoomParams, hotelId string) error
	GetRooms(ctx context.Context, hotelId string) ([]*pgtypes.Room, error)
	GetRoomById(ctx context.Context, roomId string) (*pgtypes.Room, error)
	DeleteRoom(ctx context.Context, roomId string) error
}

type PostgresRoomStore struct {
	pool *PostgresInstance
}

func NewPostgresRoomStore(pool *PostgresInstance) *PostgresRoomStore {
	return &PostgresRoomStore{
		pool: pool,
	}
}

func (s *PostgresRoomStore) CreateRoom(ctx context.Context, room pgtypes.CreateRoomParams, hotelId string) error {
	var rowId string
	row := s.pool.DB.QueryRow(ctx, `SELECT id FROM hotels WHERE id = $1`, hotelId)
	if err := row.Scan(&rowId); err != nil {
		return err
	}

	query := `INSERT INTO rooms(size, seaside, price, hotelid) VALUES($1, $2, $3, $4)`

	_, err := s.pool.DB.Exec(ctx, query, room.Size, room.SeaSide, room.Price, hotelId)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresRoomStore) GetRooms(ctx context.Context, hotelId string) ([]*pgtypes.Room, error) {
	query := `SELECT * FROM rooms WHERE hotelId = $1`

	rows, err := s.pool.DB.Query(ctx, query, hotelId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*pgtypes.Room
	for rows.Next() {
		var room pgtypes.Room
		err := rows.Scan(&room.Id, &room.Size, &room.SeaSide, &room.Price, &room.HotelId)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	return rooms, nil
}
func (s *PostgresRoomStore) GetRoomById(ctx context.Context, roomId string) (*pgtypes.Room, error) {
	query := `SELECT * FROM rooms WHERE id = $1`
	row := s.pool.DB.QueryRow(ctx, query, roomId)

	var room pgtypes.Room
	err := row.Scan(&room.Id, &room.Size, &room.SeaSide, &room.Price, &room.HotelId)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (s *PostgresRoomStore) DeleteRoom(ctx context.Context, roomId string) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := s.pool.DB.Exec(ctx, query, roomId)
	if err != nil {
		return err
	}
	return nil
}
