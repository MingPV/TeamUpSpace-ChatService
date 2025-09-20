package usecase

import (
	chatroomRepo "github.com/MingPV/ChatService/internal/chatroom/repository"
	"github.com/MingPV/ChatService/internal/entities"
	messageRepo "github.com/MingPV/ChatService/internal/message/repository"
	roommemberRepo "github.com/MingPV/ChatService/internal/room_member/repository"
)

type ChatroomService struct {
	chatroomRepository chatroomRepo.ChatroomRepository
	roommemberRepository roommemberRepo.RoomMemberRepository
	messageRepository messageRepo.MessageRepository
}

func NewChatroomService(chatroomRepository chatroomRepo.ChatroomRepository, roommemberRepository roommemberRepo.RoomMemberRepository, messageRepository messageRepo.MessageRepository) ChatroomUseCase {
	return &ChatroomService{chatroomRepository: chatroomRepository, roommemberRepository: roommemberRepository, messageRepository: messageRepository}
}

func (s *ChatroomService) CreateChatroom(chatroom *entities.Chatroom) error {
	if err := s.chatroomRepository.Save(chatroom); err != nil {
		return err
	}
	return nil
}

func (s *ChatroomService) FindChatroomByID(id int) (*entities.Chatroom, error) {
	chatroom, err := s.chatroomRepository.FindByID(id)
	if err != nil {
		return &entities.Chatroom{}, err
	}
	return chatroom, nil
}
func (s *ChatroomService) PatchChatroom(id int, chatroom *entities.Chatroom) (*entities.Chatroom, error) {
	if err := s.chatroomRepository.Patch(id, chatroom); err != nil {
		return nil, err
	}
	updatedChatroom, _ := s.chatroomRepository.FindByID(id)

	return updatedChatroom, nil
}

func (s *ChatroomService) DeleteChatroom(id int) error {
	if err := s.chatroomRepository.Delete(id); err != nil {
		return err
	}

	if err := s.messageRepository.DeleteAllMessagesByRoomID(id); err != nil {
		return err
	}

	if err := s.roommemberRepository.DeleteAllByRoomID(id); err != nil {
		return err
	}
	return nil
}