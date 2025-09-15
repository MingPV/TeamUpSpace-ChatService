package grpc

import (
	"context"

	"github.com/MingPV/ChatService/internal/chatroom/usecase"
	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/pkg/apperror"
	chatroompb "github.com/MingPV/ChatService/proto/chatroom"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcChatroomHandler struct {
	chatroomUseCase usecase.ChatroomUseCase
	chatroompb.UnimplementedChatroomServiceServer
}

func NewGrpcChatroomHandler(uc usecase.ChatroomUseCase) *GrpcChatroomHandler {
	return &GrpcChatroomHandler{chatroomUseCase: uc}
}

func (h *GrpcChatroomHandler) CreateChatroom(ctx context.Context, req *chatroompb.CreateChatroomRequest) (*chatroompb.CreateChatroomResponse, error) {
	chatroom := &entities.Chatroom{
		RoomName: req.RoomName,
		IsGroup: req.IsGroup,
	}

	if err := h.chatroomUseCase.CreateChatroom(chatroom); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &chatroompb.CreateChatroomResponse{Chatroom: toProtoChatroom(chatroom)}, nil
}

func (h *GrpcChatroomHandler) FindChatroomByID(ctx context.Context, req *chatroompb.FindChatroomByIDRequest) (*chatroompb.FindChatroomByIDResponse, error) {
	chatroom, err := h.chatroomUseCase.FindChatroomByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &chatroompb.FindChatroomByIDResponse{Chatroom: toProtoChatroom(chatroom)}, nil
}

func (h *GrpcChatroomHandler) PatchChatroom(ctx context.Context, req *chatroompb.PatchChatroomRequest) (*chatroompb.PatchChatroomResponse, error){
	chatroom := &entities.Chatroom{
		RoomName: req.RoomName,
		IsGroup: req.IsGroup,
	}
	updatedChatroom, err := h.chatroomUseCase.PatchChatroom(int(req.Id), chatroom)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &chatroompb.PatchChatroomResponse{Chatroom: toProtoChatroom(updatedChatroom)}, nil
}

func (h *GrpcChatroomHandler) DeleteChatroom(ctx context.Context, req *chatroompb.DeleteChatroomRequest) (*chatroompb.DeleteChatroomResponse, error) {
	if err := h.chatroomUseCase.DeleteChatroom(int(req.Id)); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &chatroompb.DeleteChatroomResponse{Message: "chatroom deleted"}, nil
}

func toProtoChatroom(ch *entities.Chatroom) *chatroompb.Chatroom {
	return &chatroompb.Chatroom{
		Id: int32(ch.ID),
		RoomName: ch.RoomName,
		IsGroup: ch.IsGroup,
		CreatedAt: timestamppb.New(ch.CreatedAt),
		UpdatedAt: timestamppb.New(ch.UpdatedAt),
	}
}