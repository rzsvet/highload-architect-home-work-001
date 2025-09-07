package service

import (
	"api/internal/models"
	"api/internal/monitoring"
	"api/internal/repository"
	"api/pkg/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewUserService(userRepo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) Register(user *models.User) error {
	// Check if user already exists
	exists, err := s.userRepo.UserExists(user.Username, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user with this username or email already exists")
	}

	// Validate birth date format
	if _, err := time.Parse("2006-01-02", user.BirthDate); err != nil {
		return errors.New("invalid birth date format, expected YYYY-MM-DD")
	}

	if err == nil {
		monitoring.RecordUserRegistration()
	}

	return s.userRepo.CreateUser(user)
}

func (s *UserService) Login(loginReq *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(loginReq.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(loginReq.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	monitoring.RecordUserLogin(err == nil)

	userResponse := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		Gender:    user.Gender,
		Interests: user.Interests,
		City:      user.City,
		CreatedAt: user.CreatedAt,
	}

	return &models.AuthResponse{
		Token: token,
		User:  userResponse,
	}, nil
}

func (s *UserService) GetUserByID(id int) (*models.UserResponse, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {
	return s.userRepo.GetAllUsers()
}

// func (s *UserService) UpdateUser(id int, updateReq *models.UpdateUserRequest) error {
// 	// Validate birth date if provided
// 	if updateReq.BirthDate != nil {
// 		if _, err := time.Parse("2006-01-02", *updateReq.BirthDate); err != nil {
// 			return errors.New("invalid birth date format, expected YYYY-MM-DD")
// 		}
// 	}

// 	return s.userRepo.UpdateUser(id, updateReq)
// }

func (s *UserService) generateJWT(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
