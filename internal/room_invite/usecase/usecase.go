package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/room_invite/repository"
	"github.com/google/uuid"
)

// RoomInviteService implements RoomInviteUseCase
type RoomInviteService struct {
	repo repository.RoomInviteRepository
}

// Init RoomInviteService
func NewRoomInviteService(repo repository.RoomInviteRepository) RoomInviteUseCase {
	return &RoomInviteService{repo: repo}
}

// 1. Create a new invite
func (s *RoomInviteService) CreateRoomInvite(invite *entities.RoomInvite) error {
	if err := s.repo.Save(invite); err != nil {
		return err
	}
	return nil
}

// 2. Get invite by ID
func (s *RoomInviteService) FindByID(id int) (*entities.RoomInvite, error) {
	invite, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

// 3. Get all invites by sender
func (s *RoomInviteService) FindAllBySender(sender uuid.UUID) ([]*entities.RoomInvite, error) {
	invites, err := s.repo.FindAllBySender(sender)
	if err != nil {
		return nil, err
	}
	return invites, nil
}

// 4. Get all invites by invite_to (receiver)
func (s *RoomInviteService) FindAllByInviteTo(inviteTo uuid.UUID) ([]*entities.RoomInvite, error) {
	invites, err := s.repo.FindAllByInviteTo(inviteTo)
	if err != nil {
		return nil, err
	}
	return invites, nil
}

// 5. Update invite (accept/deny etc.)
func (s *RoomInviteService) PatchInvite(id int, invite *entities.RoomInvite) error {
	if err := s.repo.Patch(id, invite); err != nil {
		return err
	}
	return nil
}

// 6. Delete invite
func (s *RoomInviteService) DeleteInvite(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
