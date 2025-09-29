package usecase

import (
	"fmt"

	chatroomRepo "github.com/MingPV/ChatService/internal/chatroom/repository"
	"github.com/MingPV/ChatService/internal/entities"
	friendRepo "github.com/MingPV/ChatService/internal/friend/repository"
	roommemberRepo "github.com/MingPV/ChatService/internal/room_member/repository"
	"github.com/google/uuid"
)


type FriendService struct {
	friendRepo friendRepo.FriendRepository
	chatroomRepo chatroomRepo.ChatroomRepository
	roommemberRepo roommemberRepo.RoomMemberRepository
}

func NewFriendService(friendRepo friendRepo.FriendRepository, chatroomRepo chatroomRepo.ChatroomRepository, roommemberRepo roommemberRepo.RoomMemberRepository) FriendUseCase {
	return &FriendService{friendRepo: friendRepo, chatroomRepo: chatroomRepo, roommemberRepo: roommemberRepo}
}

func (s *FriendService)	CreateFriend(friend *entities.Friend) error {
	if err := s.friendRepo.Save(friend); err != nil {
		return err
	}
	return nil
}

func (s *FriendService)	FindAllFriends() ([]*entities.Friend, error) {
	orders, err := s.friendRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return orders, nil
}
func (s *FriendService)	FindAllFriendsByUserID(userId uuid.UUID) ([]*entities.Friend, error) {
	orders, err := s.friendRepo.FindAllByUserId(userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *FriendService)	FindAllFriendsByIsFriend(userId uuid.UUID) ([]*entities.Friend, error) {
	orders, err := s.friendRepo.FindAllByIsFriend(userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *FriendService)	FindAllFriendsRequests(userId uuid.UUID) ([]*entities.Friend, error) {
	orders, err := s.friendRepo.FindAllFriendRequests(userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *FriendService) IsMyfriend(userId uuid.UUID, friendId uuid.UUID) (string, error) {
	status, err := s.friendRepo.IsMyfriend(userId, friendId)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (s *FriendService)	DeleteFriend(id uint) error {
	if err := s.friendRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *FriendService) FindFriendByID(id int) (*entities.Friend, error) {
	friend, err := s.friendRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func (s *FriendService) AcceptFriend(id int) (*entities.Friend, error) {
	//add friend
	friend, err := s.friendRepo.Update(id)
	if err != nil {
		return nil, err
	}

	ch, err := s.friendRepo.FindByID(id)

	//create chatroom between the two users
	chatroom := &entities.Chatroom{
		RoomName: fmt.Sprintf("room_%s_%s", ch.UserID, ch.FriendID),
		IsGroup: false,
	}
	
	if err := s.chatroomRepo.Save(chatroom); err != nil {
		return nil, err
	}

	userIDs := []uuid.UUID{ch.UserID, ch.FriendID}
	if err := s.roommemberRepo.Save(chatroom.ID, userIDs); err != nil {
		return nil, err
	}

	//create room member

	return friend, nil
}