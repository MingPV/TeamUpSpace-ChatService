package grpc

import (
	"context"

	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/room_member/usecase"
	"github.com/MingPV/ChatService/pkg/apperror"
	roommemberpb "github.com/MingPV/ChatService/proto/room_member"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcRoomMemberHandler struct {
	roomMemberUseCase usecase.RoomMemberUseCase
	roommemberpb.UnimplementedRoomMemberServiceServer
}

func NewGrpcRoomMemberHandler(uc usecase.RoomMemberUseCase) *GrpcRoomMemberHandler {
	return &GrpcRoomMemberHandler{roomMemberUseCase: uc}
}

func (h *GrpcRoomMemberHandler) CreateRoomMembers(ctx context.Context, req *roommemberpb.CreateRoomMembersRequest) (*roommemberpb.CreateRoomMembersResponse, error) {
	if err := h.roomMemberUseCase.CreateRoomMembers(uint(req.RoomId), toUUIDs(req.UserIds)); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	// after creating, fetch all members
	members, err := h.roomMemberUseCase.FindAllByRoomID(uint(req.RoomId))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoMembers []*roommemberpb.RoomMember
	for _, m := range members {
		protoMembers = append(protoMembers, toProtoRoomMember(m))
	}

	return &roommemberpb.CreateRoomMembersResponse{Members: protoMembers}, nil
}

func (h *GrpcRoomMemberHandler) FindAllByRoomID(ctx context.Context, req *roommemberpb.FindAllByRoomIDRequest) (*roommemberpb.FindAllByRoomIDResponse, error) {
	members, err := h.roomMemberUseCase.FindAllByRoomID(uint(req.RoomId))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoMembers []*roommemberpb.RoomMember
	for _, m := range members {
		protoMembers = append(protoMembers, toProtoRoomMember(m))
	}

	return &roommemberpb.FindAllByRoomIDResponse{Members: protoMembers}, nil
}

func (h *GrpcRoomMemberHandler) FindAllByUserID(ctx context.Context, req *roommemberpb.FindAllByUserIDRequest) (*roommemberpb.FindAllByUserIDResponse, error) {
	chatrooms, err := h.roomMemberUseCase.FindAllByUserID(toUUID(req.UserId))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoChatrooms []*roommemberpb.RoomMember
	for _, m := range chatrooms {
		protoChatrooms = append(protoChatrooms, toProtoRoomMember(m))
	}

	return &roommemberpb.FindAllByUserIDResponse{Chatrooms: protoChatrooms}, nil
}

func (h *GrpcRoomMemberHandler) FindByRoomIDAndUserID(ctx context.Context, req *roommemberpb.FindByRoomIDAndUserIDRequest) (*roommemberpb.FindByRoomIDAndUserIDResponse, error) {
	member, err := h.roomMemberUseCase.FindByRoomIDAndUserID(uint(req.RoomId), toUUID(req.UserId))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &roommemberpb.FindByRoomIDAndUserIDResponse{Member: toProtoRoomMember(member)}, nil
}

func (h *GrpcRoomMemberHandler) DeleteByRoomIDAndUserID(ctx context.Context, req *roommemberpb.DeleteByRoomIDAndUserIDRequest) (*roommemberpb.DeleteByRoomIDAndUserIDResponse, error) {
	if err := h.roomMemberUseCase.DeleteByRoomIDAndUserID(uint(req.RoomId), toUUID(req.UserId)); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &roommemberpb.DeleteByRoomIDAndUserIDResponse{Message: "deleted user from chatroom"}, nil
}

func (h *GrpcRoomMemberHandler) DeleteAllByRoomID(ctx context.Context, req *roommemberpb.DeleteAllByRoomIDRequest) (*roommemberpb.DeleteAllByRoomIDResponse, error) {
	if err := h.roomMemberUseCase.DeleteAllByRoomID(uint(req.RoomId)); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &roommemberpb.DeleteAllByRoomIDResponse{Message: "deleted chatroom"}, nil
}

// ---- Helper functions ----

func toProtoRoomMember(m *entities.RoomMember) *roommemberpb.RoomMember {
	return &roommemberpb.RoomMember{
		Id:        int32(m.ID),
		RoomId:    int32(m.RoomId),
		UserId:    m.UserId.String(),
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}



func toUUIDs(ids []string) []uuid.UUID {
	var uuids []uuid.UUID
	for _, id := range ids {
		if u, err := uuid.Parse(id); err == nil {
			uuids = append(uuids, u)
		}
	}
	return uuids
}

func toUUID(id string) uuid.UUID {
	u, _ := uuid.Parse(id)
	return u
}