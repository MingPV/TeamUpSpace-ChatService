package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
)

type LastVisitUseCase interface {
	// CreateLastvisit(lastvisit *entities.Lastvisit) error
	UpdateLastvisit(userId uuid.UUID, roomId int) (*entities.Lastvisit, error)
	FindByUserID(userId uuid.UUID, roomId int) (*entities.Lastvisit, error)
}