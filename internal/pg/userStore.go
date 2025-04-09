package models

import (
	"context"
	"log"

	"github.com/ctchen222/hotel-system/internal/pgtypes"
)

type PgUserStore interface {
	GetUsers(context.Context) ([]*pgtypes.PGUser, error)
	GetUserById(ctx context.Context, id string) (*pgtypes.PGUser, error)
	GetUserByEmail(ctx context.Context, email string) (*pgtypes.PGUser, error)
	CreateUser(ctx context.Context, user *pgtypes.PGUser) error
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user *pgtypes.UpdateUserParams, id string) error
}

type PostgresUserStore struct {
	pool *PostgresInstance
}

func NewPostgresUserStore(pool *PostgresInstance) *PostgresUserStore {
	return &PostgresUserStore{
		pool: pool,
	}
}

func (s *PostgresUserStore) GetUsers(ctx context.Context) ([]*pgtypes.PGUser, error) {
	query := `SELECT id, firstname, lastname, email  FROM users`

	rows, err := s.pool.DB.Query(ctx, query)
	if err != nil {
		log.Printf("Error querying users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*pgtypes.PGUser
	for rows.Next() {
		var user pgtypes.PGUser
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email); err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (s *PostgresUserStore) GetUserById(ctx context.Context, id string) (*pgtypes.PGUser, error) {
	query := `SELECT id, firstname, lastname, email FROM users WHERE id = $1`

	row := s.pool.DB.QueryRow(ctx, query, id)

	var user pgtypes.PGUser
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email); err != nil {
		log.Printf("Error scanning user: %v", err)
		return nil, err
	}

	return &user, nil
}

func (s *PostgresUserStore) GetUserByEmail(ctx context.Context, email string) (*pgtypes.PGUser, error) {
	query := `SELECT id, firstname, lastname, email, encrypted_password FROM users WHERE email = $1`

	var user pgtypes.PGUser
	row := s.pool.DB.QueryRow(ctx, query, email)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.EncryptedPassword); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *PostgresUserStore) CreateUser(ctx context.Context, user *pgtypes.PGUser) error {
	query := `INSERT INTO users(firstname, lastname, email, encrypted_password) VALUES($1, $2, $3, $4) RETURNING id, firstname, lastname`

	row := s.pool.DB.QueryRow(ctx, query, user.FirstName, user.LastName, user.Email, user.EncryptedPassword)

	var user_id, firstname, lastname string
	if err := row.Scan(&user_id, &firstname, &lastname); err != nil {
		log.Printf("Error scannind user_id: %v", err)
		return err
	}

	return nil
}

func (s *PostgresUserStore) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := s.pool.DB.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}

	return nil
}

func (s *PostgresUserStore) UpdateUser(ctx context.Context, user *pgtypes.UpdateUserParams, id string) error {
	var user_id string
	row := s.pool.DB.QueryRow(ctx, `SELECT id FROM users WHERE id = $1`, id)
	if err := row.Scan(&user_id); err != nil {
		return err
	}

	query := `UPDATE users SET firstname = $1, lastname = $2  where id = $3`

	_, err := s.pool.DB.Exec(ctx, query, user.FirstName, user.LastName, id)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}

	return nil
}
