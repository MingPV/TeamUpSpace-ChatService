package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
)

type FriendUseCase interface {
	CreateFriend(friend *entities.Friend) error 
	FindAllFriends() ([]*entities.Friend, error)
	FindAllFriendsByUserID(userId uuid.UUID) ([]*entities.Friend, error)
	FindAllFriendsByIsFriend(userId uuid.UUID) ([]*entities.Friend, error)
	FindAllFriendsRequests(userId uuid.UUID) ([]*entities.Friend, error)
	FindFriendByID(id int) (*entities.Friend, error)
	IsMyfriend(userId uuid.UUID, friendId uuid.UUID) (string, error)
	DeleteFriend(id uint) error
	AcceptFriend(id int) (*entities.Friend, error)
}