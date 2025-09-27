package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
)

type RoomInviteUseCase interface {
	CreateRoomInvite(invite *entities.RoomInvite) error
	FindByID(id int) (*entities.RoomInvite, error)
	FindAllBySender(sender uuid.UUID) ([]*entities.RoomInvite, error)
	FindAllByInviteTo(inviteTo uuid.UUID) ([]*entities.RoomInvite, error)
	FindAllByRoomId(roomId int) ([]*entities.RoomInvite, error)
	PatchInvite(id int, invite *entities.RoomInvite) error
	DeleteInvite(id int) error
	AcceptedInvite(id int) error
}
