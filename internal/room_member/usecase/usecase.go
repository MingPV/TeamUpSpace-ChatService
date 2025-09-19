package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/room_member/repository"
	"github.com/google/uuid"
)

// RoomMemberService implements RoomMemberUseCase
type RoomMemberService struct {
	repo repository.RoomMemberRepository
}

// Init RoomMemberService
func NewRoomMemberService(repo repository.RoomMemberRepository) RoomMemberUseCase {
	return &RoomMemberService{repo: repo}
}

// 1. Create multiple members in a room
func (s *RoomMemberService) CreateRoomMembers(roomId uint, userIDs []uuid.UUID) error {
	if err := s.repo.Save(roomId, userIDs); err != nil {
		return err
	}
	return nil
}

// 2. Get all members in a room
func (s *RoomMemberService) FindAllByRoomID(roomId uint) ([]*entities.RoomMember, error) {
	members, err := s.repo.FindAllByRoomID(roomId)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (s *RoomMemberService) FindAllByUserID(userId uuid.UUID) ([]*entities.RoomMember, error) {
	chatrooms, err := s.repo.FindAllByUserID(userId)
	if err != nil {
		return nil, err
	}
	return chatrooms, nil
}

// 3. Get a specific member in a room
func (s *RoomMemberService) FindByRoomIDAndUserID(roomId uint, userId uuid.UUID) (*entities.RoomMember, error) {
	member, err := s.repo.FindAllByRoomIDAndUserID(roomId, userId)
	if err != nil {
		return nil, err
	}
	return member, nil
}

// 4. Delete a specific member from a room
func (s *RoomMemberService) DeleteByRoomIDAndUserID(roomId uint, userId uuid.UUID) error {
	if err := s.repo.DeleteByRoomIDAndUserID(roomId, userId); err != nil {
		return err
	}
	return nil
}

// 5. Delete all members from a room
func (s *RoomMemberService) DeleteAllByRoomID(roomId int) error {
	if err := s.repo.DeleteAllByRoomID(roomId); err != nil {
		return err
	}
	return nil
}
