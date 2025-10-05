package grpc

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/MingPV/ChatService/pkg/apperror"

	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/message/usecase"
	messagepb "github.com/MingPV/ChatService/proto/message"
)

type GrpcMessageHandler struct {
    messageUseCase usecase.MessageUseCase
    messagepb.UnimplementedMessageServiceServer
}

func NewGrpcMessageHandler(uc usecase.MessageUseCase) *GrpcMessageHandler {
    return &GrpcMessageHandler{messageUseCase: uc}
}

func (h *GrpcMessageHandler) Chat(stream messagepb.MessageService_ChatServer) error {
    // Receive initial JoinRoom to know which room to subscribe
    var (
        roomID int
        recvErr error
    )

    // Wait for a JoinRoom event first
    first, err := stream.Recv()
    if err != nil {
        if err == io.EOF {
            return nil
        }
        return err
    }
    if join := first.GetJoin(); join != nil {
        roomID = int(join.GetRoomId())
    }
    if roomID == 0 {
        return stream.Send(&messagepb.ServerEvent{Payload: &messagepb.ServerEvent_Error{Error: &messagepb.ErrorEvent{Message: "room not specified"}}})
    }

    // Subscribe to the room and forward new messages to the client
    msgCh, cleanup := h.messageUseCase.SubscribeRoom(roomID)
    defer cleanup()

    // notify connected
    _ = stream.Send(&messagepb.ServerEvent{Payload: &messagepb.ServerEvent_Ack{Ack: &messagepb.StreamAck{Message: "joined"}}})

    // Reader goroutine: consume incoming send events and persist
    done := make(chan struct{})
    go func() {
        defer close(done)
        for {
            in, e := stream.Recv()
            if e != nil {
                recvErr = e
                return
            }
            if send := in.GetSend(); send != nil {
                senderUUID, _ := uuid.Parse(send.GetSenderId())
                now := time.Now().UTC()
                m := &entities.Message{
                    RoomId:    uint(send.GetRoomId()),
                    Message:   send.GetText(),
                    Sender:    senderUUID,
                    CreatedAt: now,
                    UpdatedAt: now,
                }
                _ = h.messageUseCase.CreateMessage(m)
            }
        }
    }()

    // Writer loop: push new messages to client
    for {
        select {
        case m, ok := <-msgCh:
            if !ok {
                return nil
            }
            _ = stream.Send(&messagepb.ServerEvent{Payload: &messagepb.ServerEvent_Delivered{Delivered: &messagepb.MessageDelivered{
                Id:             uint32(m.ID),
                RoomId:         uint32(m.RoomId),
                Text:           m.Message,
                SenderId:       m.Sender.String(),
                CreatedAtUnix:  m.CreatedAt.Unix(),
            }}})
        case <-done:
            return recvErr
        }
    }
}

func (h *GrpcMessageHandler) FindAllMessageByRoomID(ctx context.Context, req *messagepb.FindAllMessageByRoomIDRequest) (*messagepb.FindAllMessageByRoomIDResponse, error) {
    messages, err := h.messageUseCase.FindAllByRoomID(int(req.RoomId))
    if err != nil {
        return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
    }

    var protoMessages []*messagepb.Message
    for _, m := range messages {
        protoMessages = append(protoMessages, toProtoMessage(m))
    }

    return &messagepb.FindAllMessageByRoomIDResponse{Message: protoMessages}, nil

}

func (h *GrpcMessageHandler) FindLatestMessageByRoomId(ctx context.Context, req *messagepb.FindLatestMessageByRoomIdRequest) (*messagepb.FIndLastestMessageByRoomIdResponse, error) {
    message, err := h.messageUseCase.FindLatestMessageByRoomId(int(req.RoomId))
    if err != nil {
        return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
    }
    return &messagepb.FIndLastestMessageByRoomIdResponse{Message: toProtoMessage(message)}, nil
}

func toProtoMessage(m *entities.Message) *messagepb.Message {
    return &messagepb.Message{
        Id: int32(m.ID),
        RoomId: int32(m.RoomId),
        Sender: m.Sender.String(),
        Message: m.Message,
        CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.CreatedAt),
    }
}


