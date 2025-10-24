package repository

import (
	"api/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type PostRepository struct {
	writeDB *sql.DB
	readDB  *sql.DB
}

func NewPostRepository(writeDB, readDB *sql.DB) *PostRepository {
	return &PostRepository{
		writeDB: writeDB,
		readDB:  readDB,
	}
}

// CreatePost создает новый пост
func (r *PostRepository) CreatePost(post *models.Post) error {
	query := `
        INSERT INTO posts (user_id, title, content, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	now := time.Now()
	err := r.writeDB.QueryRow(
		query,
		post.UserID,
		post.Title,
		post.Content,
		now,
		now,
	).Scan(&post.ID)

	return err
}

// GetPost возвращает пост по ID
func (r *PostRepository) GetPost(postID int) (*models.PostResponse, error) {
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at,
               u.username, u.email, u.first_name, u.last_name,
               u.birth_date, u.gender, u.interests, u.city, u.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = $1
    `

	var post models.PostResponse
	var user models.UserResponse

	err := r.readDB.QueryRow(query, postID).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content,
		&post.CreatedAt, &post.UpdatedAt,
		&user.Username, &user.Email, &user.FirstName, &user.LastName,
		&user.BirthDate, &user.Gender, &user.Interests, &user.City, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, err
	}

	post.User = user
	return &post, nil
}

// UpdatePost обновляет пост
func (r *PostRepository) UpdatePost(postID, userID int, updateReq *models.UpdatePostRequest) error {
	query := `
        UPDATE posts 
        SET title = COALESCE($1, title),
            content = COALESCE($2, content),
            updated_at = $3
        WHERE id = $4 AND user_id = $5
    `

	result, err := r.writeDB.Exec(
		query,
		updateReq.Title,
		updateReq.Content,
		time.Now(),
		postID,
		userID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("post not found or access denied")
	}

	return nil
}

// DeletePost удаляет пост
func (r *PostRepository) DeletePost(postID, userID int) error {
	query := `DELETE FROM posts WHERE id = $1 AND user_id = $2`

	result, err := r.writeDB.Exec(query, postID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("post not found or access denied")
	}

	return nil
}

// GetUserPosts возвращает посты пользователя
func (r *PostRepository) GetUserPosts(userID, limit, offset int) ([]models.PostResponse, int, error) {
	// Счетчик общего количества
	var total int
	countQuery := `SELECT COUNT(*) FROM posts WHERE user_id = $1`
	err := r.readDB.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Получение постов
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at,
               u.username, u.email, u.first_name, u.last_name,
               u.birth_date, u.gender, u.interests, u.city, u.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.user_id = $1
        ORDER BY p.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.readDB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []models.PostResponse
	for rows.Next() {
		var post models.PostResponse
		var user models.UserResponse

		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content,
			&post.CreatedAt, &post.UpdatedAt,
			&user.Username, &user.Email, &user.FirstName, &user.LastName,
			&user.BirthDate, &user.Gender, &user.Interests, &user.City, &user.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		post.User = user
		posts = append(posts, post)
	}

	return posts, total, nil
}

// GetFriendsPosts возвращает посты друзей пользователя (лента)
func (r *PostRepository) GetFriendsPosts(userID, limit, offset int) ([]models.PostResponse, int, error) {
	// Счетчик общего количества
	var total int
	countQuery := `
        SELECT COUNT(*) 
        FROM posts p
        JOIN friends f ON p.user_id = f.friend_id
        WHERE f.user_id = $1
    `
	err := r.readDB.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Получение постов друзей
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at,
               u.username, u.email, u.first_name, u.last_name,
               u.birth_date, u.gender, u.interests, u.city, u.created_at
        FROM posts p
        JOIN friends f ON p.user_id = f.friend_id
        JOIN users u ON p.user_id = u.id
        WHERE f.user_id = $1
        ORDER BY p.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.readDB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []models.PostResponse
	for rows.Next() {
		var post models.PostResponse
		var user models.UserResponse

		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content,
			&post.CreatedAt, &post.UpdatedAt,
			&user.Username, &user.Email, &user.FirstName, &user.LastName,
			&user.BirthDate, &user.Gender, &user.Interests, &user.City, &user.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		post.User = user
		posts = append(posts, post)
	}

	return posts, total, nil
}
