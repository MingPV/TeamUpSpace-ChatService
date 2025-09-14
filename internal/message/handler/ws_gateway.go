package handler

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	messagepb "github.com/MingPV/ChatService/proto/message"
)

type WebSocketGatewayHandler struct {
    client messagepb.MessageServiceClient
}

func NewWebSocketGatewayHandler(client messagepb.MessageServiceClient) *WebSocketGatewayHandler {
    return &WebSocketGatewayHandler{client: client}
}

// UpgradeMiddleware ensures the request is a valid WebSocket upgrade
func (h *WebSocketGatewayHandler) UpgradeMiddleware(c *fiber.Ctx) error {
    if websocket.IsWebSocketUpgrade(c) {
        return c.Next()
    }
    return fiber.ErrUpgradeRequired
}

// SubscribeRoomWebSocket bridges the WebSocket to the gRPC streaming Chat method.
func (h *WebSocketGatewayHandler) SubscribeRoomWebSocket(c *websocket.Conn) {
    roomIDStr := c.Params("roomId")
    roomID, err := strconv.Atoi(roomIDStr)
    if err != nil || roomID <= 0 {
        _ = c.WriteMessage(websocket.TextMessage, []byte("invalid room id"))
        _ = c.Close()
        return
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    stream, err := h.client.Chat(ctx)
    if err != nil {
        _ = c.WriteMessage(websocket.TextMessage, []byte("failed to connect stream"))
        _ = c.Close()
        return
    }

    // Send JoinRoom first
    _ = stream.Send(&messagepb.ClientEvent{Payload: &messagepb.ClientEvent_Join{Join: &messagepb.JoinRoom{RoomId: uint32(roomID)}}})

    // Reader: WS -> gRPC stream
    type inbound struct {
        Message string `json:"message"`
        Sender  string `json:"sender"`
        SentAt  int64  `json:"sent_at_unix"`
    }
    go func() {
        for {
            _, data, rerr := c.ReadMessage()
            if rerr != nil {
                _ = stream.CloseSend()
                return
            }
            var in inbound
            if jerr := json.Unmarshal(data, &in); jerr != nil {
                continue
            }
            sentAt := in.SentAt
            if sentAt == 0 {
                sentAt = time.Now().Unix()
            }
            _ = stream.Send(&messagepb.ClientEvent{Payload: &messagepb.ClientEvent_Send{Send: &messagepb.SendMessage{
                RoomId:      uint32(roomID),
                Text:        in.Message,
                SenderId:    in.Sender,
                SentAtUnix:  sentAt,
            }}})
        }
    }()

    // Writer: gRPC stream -> WS
    for {
        ev, rerr := stream.Recv()
        if rerr != nil {
            _ = c.WriteMessage(websocket.TextMessage, []byte("stream closed"))
            return
        }
        // Forward server event as JSON
        b, _ := json.Marshal(ev)
        _ = c.WriteMessage(websocket.TextMessage, b)
    }
}


