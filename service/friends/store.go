package friends

import (
	"database/sql"
	"fmt"
	"github.com/Code-xCartel/noxus-api-svc/service/auth"
	"github.com/Code-xCartel/noxus-api-svc/types/friends"
	"github.com/Code-xCartel/noxus-api-svc/utils"
	"github.com/lib/pq"
)

type Store struct {
	db        *sql.DB
	authStore *auth.Store
}

type Status string

const (
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Rejected Status = "rejected"
	Blocked  Status = "blocked"
)

func NewFriendsStore(db *sql.DB, authStore *auth.Store) *Store {
	return &Store{db, authStore}
}

func (s *Store) GetFriends(id string, status Status) ([]friends.FriendResponse, error) {
	query := `
		SELECT u.noxId, u.username, f.status 
		FROM users u 
		JOIN friends f 
		ON (u.noxId = f.friend_id OR u.noxId = f.user_id) 
		WHERE (f.user_id = $1 OR f.friend_id = $1) AND  f.status = $2
	`

	if status == Blocked {
		query += `AND f.action_by = $1`
	}

	rows, dbErr := s.db.Query(query, id, status)
	if dbErr != nil {
		return nil, dbErr
	}
	var users []friends.FriendResponse
	for rows.Next() {
		u, err := scanRowIntoModel(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}

	filterSelf := func(user friends.FriendResponse) bool { return user.NoxID != id }
	filteredUsers := utils.FilterArray(users, filterSelf)

	return filteredUsers, nil
}

func (s *Store) CheckFriendStatusByNoxId(userId string, friendId string, status Status) (bool, error) {
	rows, err := s.db.Query(
		"SELECT u.noxId, u.username, f.status FROM users u JOIN friends f ON (u.noxId = f.user_id) WHERE (f.user_id = $1 AND f.friend_id = $2) AND  f.status = $3",
		userId, friendId, status)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		return false, fmt.Errorf("relation not found")
	}
	return true, nil
}

func (s *Store) AddFriendByNoxId(userId string, friendId string) error {
	_, err := s.db.Exec(
		"INSERT INTO friends (user_id, friend_id, status) VALUES ($1, $2, 'pending') ON CONFLICT DO NOTHING",
		userId, friendId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteFriendByNoxId(userId string, friendId string) error {
	_, err := s.db.Exec(
		"DELETE FROM friends WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)",
		userId, friendId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ActionOnFriendByNoxId(
	userId string,
	friendId string,
	currentStatus []Status,
	newStatus Status,
) error {
	_, err := s.db.Exec(
		"UPDATE friends SET status = $4, action_by = $2 WHERE (user_id = $1 AND friend_id = $2 OR user_id = $2 AND friend_id = $1) AND status = ANY($3)",
		friendId, userId, pq.Array(currentStatus), newStatus,
	)
	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoModel(row *sql.Rows) (*friends.FriendResponse, error) {
	user := new(friends.FriendResponse)
	if err := row.Scan(
		&user.NoxID,
		&user.Username,
		&user.Status,
	); err != nil {
		return nil, err
	}
	friend := &friends.FriendResponse{
		Username: user.Username,
		NoxID:    user.NoxID,
		Status:   user.Status,
	}
	return friend, nil
}
