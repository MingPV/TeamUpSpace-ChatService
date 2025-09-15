package repository

import (
	"github.com/MingPV/ChatService/internal/entities"
)

type ChatroomRepository interface {
	Save(chatroom *entities.Chatroom) error 
	Patch(id int, chatroom *entities.Chatroom) error
	FindByID(id int) (*entities.Chatroom, error)
	Delete(id int) error
}