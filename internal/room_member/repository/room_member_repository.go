package repository

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
)

type RoomMemberRepository interface {
	Save(roomId uint, userIDs []uuid.UUID) error
	FindAllByRoomID(roomID uint) ([]*entities.RoomMember, error)
	FindAllByUserID(userId uuid.UUID) ([]*entities.RoomMember, error)
	FindAllByRoomIDAndUserID(roomId uint, userId uuid.UUID) (*entities.RoomMember, error)
	DeleteByRoomIDAndUserID(roomID uint, userID uuid.UUID) error
	DeleteAllByRoomID(roomID int) error
	Delete(id int) error
}

