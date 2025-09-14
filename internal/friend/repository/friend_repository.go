package repository

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
)

type FriendRepository interface {
	Save(friend *entities.Friend) error
	FindAll() ([]*entities.Friend, error)
	FindAllByUserId(userId uuid.UUID) ([]*entities.Friend, error)
	FindAllByIsFriend(userId uuid.UUID) ([]*entities.Friend, error)
	FindAllFriendRequests(userId uuid.UUID) ([]*entities.Friend, error)
	IsMyfriend(userId uuid.UUID, friendId uuid.UUID) (string, error)
	Delete(id uint) error 
	Update(userId uuid.UUID, friendId uuid.UUID) (*entities.Friend, error)
}