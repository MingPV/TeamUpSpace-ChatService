package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
)

type MessageUseCase interface {
	CreateMessage(message *entities.Message) error
	FindAllByRoomID(roomId int) ([]*entities.Message, error)
	DeleteAllMessagesByRoomID(roomId int) error
	FindLatestMessageByRoomId(roomId int) (*entities.Message, error)
	FindAllMessagesUnread(userId uuid.UUID) ([]*entities.Message, error)

	// SubscribeRoom subscribes to a room and returns a read-only channel of messages
	// and a cleanup function to unsubscribe and release resources.
	SubscribeRoom(roomId int) (<-chan *entities.Message, func())
}