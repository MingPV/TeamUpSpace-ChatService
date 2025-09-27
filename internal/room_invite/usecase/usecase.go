package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	roominviteRepo "github.com/MingPV/ChatService/internal/room_invite/repository"
	roommemberRepo "github.com/MingPV/ChatService/internal/room_member/repository"
	"github.com/google/uuid"
)

// RoomInviteService implements RoomInviteUseCase
type RoomInviteService struct {
	roominviteRepo roominviteRepo.RoomInviteRepository
	roommemberRepo roommemberRepo.RoomMemberRepository
}

// Init RoomInviteService
func NewRoomInviteService(roominviteRepo roominviteRepo.RoomInviteRepository, roommemberRepo roommemberRepo.RoomMemberRepository) RoomInviteUseCase {
	return &RoomInviteService{roominviteRepo: roominviteRepo, roommemberRepo: roommemberRepo}
}

// 1. Create a new invite
func (s *RoomInviteService) CreateRoomInvite(invite *entities.RoomInvite) error {
	if err := s.roominviteRepo.Save(invite); err != nil {
		return err
	}
	return nil
}

// 2. Get invite by ID
func (s *RoomInviteService) FindByID(id int) (*entities.RoomInvite, error) {
	invite, err := s.roominviteRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

// 3. Get all invites by sender
func (s *RoomInviteService) FindAllBySender(sender uuid.UUID) ([]*entities.RoomInvite, error) {
	invites, err := s.roominviteRepo.FindAllBySender(sender)
	if err != nil {
		return nil, err
	}
	return invites, nil
}

// 4. Get all invites by invite_to (receiver)
func (s *RoomInviteService) FindAllByInviteTo(inviteTo uuid.UUID) ([]*entities.RoomInvite, error) {
	invites, err := s.roominviteRepo.FindAllByInviteTo(inviteTo)
	if err != nil {
		return nil, err
	}
	return invites, nil
}

func (s *RoomInviteService) FindAllByRoomId(roomId int) ([]*entities.RoomInvite, error) {
	invites, err := s.roominviteRepo.FindAllByRoomId(int(roomId))
	if err != nil {
		return nil, err
	}
	return invites, nil
}

// 5. Update invite (accept/deny etc.)
func (s *RoomInviteService) PatchInvite(id int, invite *entities.RoomInvite) error {
	if err := s.roominviteRepo.Patch(id, invite); err != nil {
		return err
	}
	return nil
}

// 6. Delete invite
func (s *RoomInviteService) DeleteInvite(id int) error {
	if err := s.roominviteRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *RoomInviteService) AcceptedInvite(id int) error {
	invite, err := s.roominviteRepo.FindByID(id)
	if err != nil {
		return err;
	}

	if err := s.roommemberRepo.Save(invite.RoomId, []uuid.UUID{invite.InviteTo}); err != nil {
        return err
    }

	if err := s.roominviteRepo.Delete(id); err != nil {
		return err
	}
	return nil
}
