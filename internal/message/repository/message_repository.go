package repository

import (
	"github.com/MingPV/ChatService/internal/entities"
)

type MessageRepository interface {
	Save(message *entities.Message) error
	FindAllByRoomID(roomId int) ([]*entities.Message, error)

	DeleteAllMessagesByRoomID(roomId int) error
	FindByRoomId(roomId int) (*entities.Message, error)

}