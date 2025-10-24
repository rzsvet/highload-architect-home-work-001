package repository

import (
	"api/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type FriendRepository struct {
	writeDB *sql.DB
	readDB  *sql.DB
}

func NewFriendRepository(writeDB, readDB *sql.DB) *FriendRepository {
	return &FriendRepository{
		writeDB: writeDB,
		readDB:  readDB,
	}
}

// AddFriend добавляет друга
func (r *FriendRepository) AddFriend(userID, friendID int) error {
	// Проверяем, что пользователь не пытается добавить сам себя
	if userID == friendID {
		return fmt.Errorf("cannot add yourself as a friend")
	}

	// Проверяем существование пользователя
	if exists, err := r.userExists(friendID); err != nil || !exists {
		return fmt.Errorf("friend user does not exist")
	}

	// Проверяем, не добавлен ли уже друг
	if isFriend, err := r.IsFriend(userID, friendID); err != nil || isFriend {
		return fmt.Errorf("users are already friends")
	}

	query := `
        INSERT INTO friends (user_id, friend_id, created_at) 
        VALUES ($1, $2, $3), ($2, $1, $3)
    `

	_, err := r.writeDB.Exec(query, userID, friendID, time.Now())
	return err
}

// DeleteFriend удаляет друга
func (r *FriendRepository) DeleteFriend(userID, friendID int) error {
	query := `
        DELETE FROM friends 
        WHERE (user_id = $1 AND friend_id = $2) 
           OR (user_id = $2 AND friend_id = $1)
    `

	result, err := r.writeDB.Exec(query, userID, friendID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("friendship not found")
	}

	return nil
}

// IsFriend проверяет, являются ли пользователи друзьями
func (r *FriendRepository) IsFriend(userID, friendID int) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM friends 
            WHERE user_id = $1 AND friend_id = $2
        )
    `

	var exists bool
	err := r.readDB.QueryRow(query, userID, friendID).Scan(&exists)
	return exists, err
}

// GetFriends возвращает список друзей пользователя
func (r *FriendRepository) GetFriends(userID int) ([]models.FriendResponse, error) {
	query := `
        SELECT f.id, f.user_id, f.friend_id, f.created_at,
               u.username, u.email, u.first_name, u.last_name,
               u.birth_date, u.gender, u.interests, u.city, u.created_at
        FROM friends f
        JOIN users u ON f.friend_id = u.id
        WHERE f.user_id = $1
        ORDER BY f.created_at DESC
    `

	rows, err := r.readDB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.FriendResponse
	for rows.Next() {
		var friend models.FriendResponse
		var user models.UserResponse

		err := rows.Scan(
			&friend.ID, &friend.UserID, &friend.FriendID, &friend.CreatedAt,
			&user.Username, &user.Email, &user.FirstName, &user.LastName,
			&user.BirthDate, &user.Gender, &user.Interests, &user.City, &user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		friend.Friend = user
		friends = append(friends, friend)
	}

	return friends, nil
}

// GetFriendshipStatus возвращает статус дружбы между пользователями
func (r *FriendRepository) GetFriendshipStatus(userID, friendID int) (*models.FriendshipStatus, error) {
	isFriend, err := r.IsFriend(userID, friendID)
	if err != nil {
		return nil, err
	}

	return &models.FriendshipStatus{
		IsFriend:  isFriend,
		IsPending: false, // В расширенной версии можно добавить pending статус
	}, nil
}

func (r *FriendRepository) userExists(userID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	err := r.readDB.QueryRow(query, userID).Scan(&exists)
	return exists, err
}
