package auth

import (
	"database/sql"
	"fmt"
	"github.com/Code-xCartel/noxus-api-svc/types/auth"
)

type Store struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) GetUserByEmail(email string) (*auth.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	u := new(auth.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) GetUserByNoxID(noxId string) (*auth.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE noxId = $1", noxId)
	if err != nil {
		return nil, err
	}
	u := new(auth.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) CreateNewUser(u auth.User) error {
	if _, err := s.db.Exec("INSERT INTO users (noxId, username, email, password) VALUES ($1, $2, $3, $4)",
		u.NoxID,
		u.Username,
		u.Email,
		u.Password,
	); err != nil {
		return err
	}
	return nil
}

func scanRowIntoUser(row *sql.Rows) (*auth.User, error) {
	user := new(auth.User)
	if err := row.Scan(
		&user.ID,
		&user.NoxID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}
	return user, nil
}
