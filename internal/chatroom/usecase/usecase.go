package usecase

import (
	"github.com/MingPV/ChatService/internal/chatroom/repository"
	"github.com/MingPV/ChatService/internal/entities"
)

type ChatroomService struct {
	repo repository.ChatroomRepository
}

func NewChatroomService(repo repository.ChatroomRepository) ChatroomUseCase {
	return &ChatroomService{repo: repo}
}

func (s *ChatroomService) CreateChatroom(chatroom *entities.Chatroom) error {
	if err := s.repo.Save(chatroom); err != nil {
		return err
	}
	return nil
}

func (s *ChatroomService) FindChatroomByID(id int) (*entities.Chatroom, error) {
	chatroom, err := s.repo.FindByID(id)
	if err != nil {
		return &entities.Chatroom{}, err
	}
	return chatroom, nil
}
func (s *ChatroomService) PatchChatroom(id int, chatroom *entities.Chatroom) (*entities.Chatroom, error) {
	if err := s.repo.Patch(id, chatroom); err != nil {
		return nil, err
	}
	updatedChatroom, _ := s.repo.FindByID(id)

	return updatedChatroom, nil
}

func (s *ChatroomService) DeleteChatroom(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}