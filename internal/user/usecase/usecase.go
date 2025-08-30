package usecase

import (
	"os"
	"time"

	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/user/repository"
	"github.com/MingPV/ChatService/pkg/apperror"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ChatService struct
type ChatService struct {
	repo repository.UserRepository
}

// Init ChatService
func NewChatService(repo repository.UserRepository) UserUseCase {
	return &ChatService{repo: repo}
}

// ChatService Methods - 1 Register user (hash password)
func (s *ChatService) Register(user *entities.User) error {
	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser != nil {
		return apperror.ErrAlreadyExists
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPwd)

	return s.repo.Save(user)
}

// ChatService Methods - 2 Login user (check email + password)
func (s *ChatService) Login(email string, password string) (string, *entities.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, err
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 3 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

// ChatService Methods - 3 Get user by id
func (s *ChatService) FindUserByID(id string) (*entities.User, error) {
	return s.repo.FindByID(id)
}

// ChatService Methods - 4 Get all users
func (s *ChatService) FindAllUsers() ([]*entities.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ChatService Methods - 5 Get user by email
func (s *ChatService) GetUserByEmail(email string) (*entities.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ChatService Methods - 6 Patch
func (s *ChatService) PatchUser(id string, user *entities.User) (*entities.User, error) {
	if err := s.repo.Patch(id, user); err != nil {
		return nil, err
	}
	updatedUser, _ := s.repo.FindByID(id)

	return updatedUser, nil
}

// ChatService Methods - 7 Delete
func (s *ChatService) DeleteUser(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
