package repository

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
)

type RoomInviteRepository interface {
	Save(invite *entities.RoomInvite) error
	FindByID(id int) (*entities.RoomInvite, error)
	FindAllBySender(sender uuid.UUID) ([]*entities.RoomInvite, error)
	FindAllByRoomId(roomId int) ([]*entities.RoomInvite, error)
	FindAllByInviteTo(inviteTo uuid.UUID) ([]*entities.RoomInvite, error)
	Patch(id int, invite *entities.RoomInvite) error
	Delete(id int) error
}