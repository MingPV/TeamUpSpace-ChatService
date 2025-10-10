package usecase

import (
	"sync"

	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/message/repository"
	"github.com/google/uuid"
)

type MessageService struct {
	repo repository.MessageRepository

	subscribers map[int][]chan *entities.Message // roomId -> list of channels
	mu          sync.RWMutex
}

func NewMessageService(repo repository.MessageRepository) MessageUseCase {
	return &MessageService{repo: repo, subscribers: make(map[int][]chan *entities.Message),}
}

func (s *MessageService) CreateMessage(message *entities.Message) error {
	if err := s.repo.Save(message); err != nil {
		return err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ch := range s.subscribers[int(message.RoomId)] {
		select {
		case ch <- message: // non-blocking
		default:
		}
	}
	return nil
}

func (s *MessageService) FindAllByRoomID(roomId int) ([]*entities.Message, error) {
	messages, err := s.repo.FindAllByRoomID(roomId)
	if err != nil {
		return nil, err
	}
	return messages, nil
} 

func (s *MessageService) SubscribeRoom(roomId int) (<-chan *entities.Message, func()) {
	ch := make(chan *entities.Message, 10) // buffered channel

	// Add channel to subscribers
	s.mu.Lock()
	s.subscribers[roomId] = append(s.subscribers[roomId], ch)
	s.mu.Unlock()

	// Return the channel and a cleanup function
	cleanup := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		channels := s.subscribers[roomId]
		for i, c := range channels {
			if c == ch {
				// remove from slice
				s.subscribers[roomId] = append(channels[:i], channels[i+1:]...)
				close(c)
				break
			}
		}
	}

	return ch, cleanup
}

func (s *MessageService) DeleteAllMessagesByRoomID(roomId int) error {
	if err := s.repo.DeleteAllMessagesByRoomID(roomId); err != nil {
		return err
	}
	return nil
}

func (s *MessageService) FindLatestMessageByRoomId(roomId int) (*entities.Message, error) {
	message, err := s.repo.FindByRoomId(roomId); 
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MessageService) FindAllMessagesUnread(userId uuid.UUID) ([]*entities.Message, error) {
	messages, err := s.repo.FindAllMessagesUnread(userId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}