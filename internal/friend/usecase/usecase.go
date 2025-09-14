package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/friend/repository"
	"github.com/google/uuid"
)


type FriendService struct {
	repo repository.FriendRepository
}

func NewFriendService(repo repository.FriendRepository) FriendUseCase {
	return &FriendService{repo: repo}
}

func (s *FriendService)	CreateFriend(friend *entities.Friend) error {
	if err := s.repo.Save(friend); err != nil {
		return err
	}
	return nil
}

func (s *FriendService)	FindAllFriends() ([]*entities.Friend, error) {
	orders, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return orders, nil
}
func (s *FriendService)	FindAllFriendsByUserID(userId uuid.UUID) ([]*entities.Friend, error) {
	orders, err := s.repo.FindAllByUserId(userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *FriendService)	FindAllFriendsByIsFriend(userId uuid.UUID) ([]*entities.Friend, error) {
	orders, err := s.repo.FindAllByIsFriend(userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *FriendService)	FindAllFriendsRequests(userId uuid.UUID) ([]*entities.Friend, error) {
	orders, err := s.repo.FindAllFriendRequests(userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *FriendService) IsMyfriend(userId uuid.UUID, friendId uuid.UUID) (string, error) {
	status, err := s.repo.IsMyfriend(userId, friendId)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (s *FriendService)	DeleteFriend(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *FriendService) AcceptFriend(userId uuid.UUID, friendId uuid.UUID) (*entities.Friend, error) {
	friend, err := s.repo.Update(userId, friendId)
	if err != nil {
		return nil, err
	}
	return friend, nil
}