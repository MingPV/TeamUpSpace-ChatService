package usecase

import (
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/lastvisit/repository"
	"github.com/google/uuid"
)

type LastvisitService struct {
	repo repository.LastvisitRepository
}

func NewLastvisitService(repo repository.LastvisitRepository) LastVisitUseCase {
	return &LastvisitService{repo: repo}
}


// func (s *LastvisitService) CreateLastvisit(lastvisit *entities.Lastvisit) error {
// 	if err := s.repo.Save(lastvisit); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (s *LastvisitService) UpdateLastvisit(userId uuid.UUID, roomId int) (*entities.Lastvisit, error) {
    updatedLastvisit, err := s.repo.Patch(userId, roomId)
    if err != nil {
        return nil, err
    }

    return updatedLastvisit, nil
}


func (s *LastvisitService) FindByUserID(userId uuid.UUID, roomId int) (*entities.Lastvisit, error){
	lastvisit, err := s.repo.FindByUserId(userId, roomId)
	if err != nil {
		return nil, err
	}
	return lastvisit, nil
}