package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
)

type RoomMemberUseCase interface {
	CreateRoomMembers(roomId uint, userIDs []uuid.UUID) error
	FindAllByRoomID(roomId uint) ([]*entities.RoomMember, error)
	FindAllByUserID(userId uuid.UUID) ([]*entities.RoomMember, error)
	FindByRoomIDAndUserID(roomId uint, userId uuid.UUID) (*entities.RoomMember, error)
	DeleteByRoomIDAndUserID(roomId uint, userId uuid.UUID) error
	DeleteAllByRoomID(roomId int) error
	DeleteRoomMember(id int) error
}
