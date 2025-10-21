package repository

import (
	"api/internal/models"
	"api/internal/monitoring"
	"api/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type UserRepository struct {
	writeDB *sql.DB
	readDB  *sql.DB
}

func NewUserRepository(writeDB, readDB *sql.DB) *UserRepository {
	return &UserRepository{
		writeDB: writeDB,
		readDB:  readDB,
	}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	start := time.Now()

	query := `
        INSERT INTO users (
            username, email, password, first_name, last_name, 
            birth_date, gender, interests, city, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id
    `

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	err = r.writeDB.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		hashedPassword,
		user.FirstName,
		user.LastName,
		user.BirthDate,
		user.Gender,
		user.Interests,
		user.City,
		now,
		now,
	).Scan(&user.ID)

	defer func() {
		duration := time.Since(start)
		monitoring.ObserveDatabaseQuery("create_user", err == nil, duration)
	}()

	return err
}

func (r *UserRepository) GetUserByID(id int) (*models.UserResponse, error) {
	start := time.Now()

	query := `
        SELECT 
            id, username, email, first_name, last_name, 
            birth_date, gender, interests, city, created_at
        FROM users 
        WHERE id = $1
    `

	var user models.UserResponse
	err := r.readDB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.BirthDate,
		&user.Gender,
		&user.Interests,
		&user.City,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		duration := time.Since(start)
		monitoring.ObserveDatabaseQuery("get_user_by_id", err == nil, duration)
	}()

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	start := time.Now()

	query := `
        SELECT
            id, username, email, password, first_name, last_name,
            birth_date, gender, interests, city, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var user models.User
	err := r.readDB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.BirthDate,
		&user.Gender,
		&user.Interests,
		&user.City,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	defer func() {
		duration := time.Since(start)
		monitoring.ObserveDatabaseQuery("get_all_users", err == nil, duration)
	}()

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]models.UserResponse, error) {
	start := time.Now()

	query := `
        SELECT 
            id, username, email, first_name, last_name, 
            birth_date, gender, interests, city, created_at
        FROM users 
        ORDER BY created_at DESC
    `

	rows, err := r.readDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	defer func() {
		duration := time.Since(start)
		monitoring.ObserveDatabaseQuery("get_all_users", err == nil, duration)
	}()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.BirthDate,
			&user.Gender,
			&user.Interests,
			&user.City,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) UserExists(username string, email string) (bool, error) {
	start := time.Now()
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)`
	var exists bool
	err := r.readDB.QueryRow(query, username, email).Scan(&exists)
	defer func() {
		duration := time.Since(start)
		monitoring.ObserveDatabaseQuery("user_exists", err == nil, duration)
	}()
	return exists, err
}

func (r *UserRepository) SearchUsers(firstName, lastName string) ([]models.UserResponse, error) {
	query := `
        SELECT 
            id, username, email, first_name, last_name, 
            birth_date, gender, interests, city, created_at
        FROM users 
        WHERE first_name ILIKE $1 AND last_name ILIKE $2
        ORDER BY id
        LIMIT 100
    `

	// query := `
	//     SELECT
	//         id, username, email, first_name, last_name,
	//         birth_date, gender, interests, city, created_at
	//     FROM users
	//     WHERE first_name ILIKE $1 AND last_name ILIKE $2
	//     ORDER BY first_name, last_name
	//     LIMIT 100
	// `

	// Добавляем wildcards для частичного совпадения
	firstNamePattern := "%" + strings.ToLower(firstName) + "%"
	lastNamePattern := "%" + strings.ToLower(lastName) + "%"

	rows, err := r.readDB.Query(query, firstNamePattern, lastNamePattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.BirthDate,
			&user.Gender,
			&user.Interests,
			&user.City,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// SearchUsersWithPaging поиск с пагинацией
func (r *UserRepository) SearchUsersWithPaging(firstName string, lastName string, limit int, offset int) ([]models.UserResponse, int, error) {
	// Запрос для получения пользователей
	usersQuery := `
        SELECT 
            id, username, email, first_name, last_name, 
            birth_date, gender, interests, city, created_at
        FROM users 
        WHERE first_name ILIKE $1 AND last_name ILIKE $2
        ORDER BY id
        LIMIT $3 OFFSET $4
    `

	// usersQuery := `
	//     SELECT
	//         id, username, email, first_name, last_name,
	//         birth_date, gender, interests, city, created_at
	//     FROM users
	//     WHERE first_name ILIKE $1 AND last_name ILIKE $2
	//     ORDER BY first_name, last_name
	//     LIMIT $3 OFFSET $4
	// `

	// Запрос для подсчета общего количества
	countQuery := `
        SELECT COUNT(*) 
        FROM users 
        WHERE first_name ILIKE $1 AND last_name ILIKE $2
    `

	firstNamePattern := "%" + strings.ToLower(firstName) + "%"
	lastNamePattern := "%" + strings.ToLower(lastName) + "%"

	// Выполняем запрос на подсчет
	var total int
	err := r.readDB.QueryRow(countQuery, firstNamePattern, lastNamePattern).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Выполняем запрос на получение данных
	rows, err := r.readDB.Query(usersQuery, firstNamePattern, lastNamePattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.BirthDate,
			&user.Gender,
			&user.Interests,
			&user.City,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// func (r *UserRepository) UpdateUser(id int, updateReq *models.UpdateUserRequest) error {
// 	query := `
//         UPDATE users
//         SET
//             first_name = COALESCE($1, first_name),
//             last_name = COALESCE($2, last_name),
//             birth_date = COALESCE($3, birth_date),
//             gender = COALESCE($4, gender),
//             interests = COALESCE($5, interests),
//             city = COALESCE($6, city),
//             updated_at = $7
//         WHERE id = $8
//     `

// 	_, err := r.writeDB.Exec(
// 		query,
// 		updateReq.FirstName,
// 		updateReq.LastName,
// 		updateReq.BirthDate,
// 		updateReq.Gender,
// 		updateReq.Interests,
// 		updateReq.City,
// 		time.Now(),
// 		id,
// 	)

// 	return err
// }
