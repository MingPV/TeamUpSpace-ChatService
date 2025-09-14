package handler

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"

	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/message/usecase"
)

type WebSocketMessageHandler struct {
    messageUseCase usecase.MessageUseCase
}

func NewWebSocketMessageHandler(uc usecase.MessageUseCase) *WebSocketMessageHandler {
    return &WebSocketMessageHandler{messageUseCase: uc}
}

// UpgradeMiddleware ensures the request is a valid WebSocket upgrade
func (h *WebSocketMessageHandler) UpgradeMiddleware(c *fiber.Ctx) error {
    if websocket.IsWebSocketUpgrade(c) {
        return c.Next()
    }
    return fiber.ErrUpgradeRequired
}

// SubscribeRoomWebSocket handles WebSocket connections for a specific room.
// Path params: :roomId
// Optional query params: sender=<uuid>
func (h *WebSocketMessageHandler) SubscribeRoomWebSocket(c *websocket.Conn) {
    roomIDStr := c.Params("roomId")
    roomID, err := strconv.Atoi(roomIDStr)
    if err != nil || roomID <= 0 {
        _ = c.WriteMessage(websocket.TextMessage, []byte("invalid room id"))
        _ = c.Close()
        return
    }

    // Subscribe to the room
    msgCh, cleanup := h.messageUseCase.SubscribeRoom(roomID)
    defer cleanup()

    // Reader goroutine: receive messages from client and persist
    type inbound struct {
        Message string    `json:"message"`
        Sender  string    `json:"sender"`
        SentAt  time.Time `json:"sent_at"`
    }

    go func() {
        for {
            _, data, rerr := c.ReadMessage()
            if rerr != nil {
                return
            }
            var in inbound
            if jerr := json.Unmarshal(data, &in); jerr != nil {
                continue
            }
            senderUUID, uerr := uuid.Parse(in.Sender)
            if uerr != nil {
                continue
            }
            now := time.Now().UTC()
            m := &entities.Message{
                RoomId:    uint(roomID),
                Message:   in.Message,
                Sender:    senderUUID,
                CreatedAt: now,
                UpdatedAt: now,
            }
            _ = h.messageUseCase.CreateMessage(m)
        }
    }()

    // Writer loop: forward new messages to this client
    for m := range msgCh {
        _ = c.WriteJSON(m)
    }
}


