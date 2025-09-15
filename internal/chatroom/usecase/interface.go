package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
)

type ChatroomUseCase interface {
	CreateChatroom(chatroom *entities.Chatroom) error 
	FindChatroomByID(id int) (*entities.Chatroom, error)
	PatchChatroom(id int, chatroom *entities.Chatroom) (*entities.Chatroom, error)
	DeleteChatroom(id int) error
}