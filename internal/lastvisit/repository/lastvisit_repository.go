package repository

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
)

type LastvisitRepository interface {
	// Save(lastvisit *entities.Lastvisit) error
	FindByUserId(userId uuid.UUID, roomId int) (*entities.Lastvisit, error)
	Patch(userId uuid.UUID, roomId int) (*entities.Lastvisit, error)
}