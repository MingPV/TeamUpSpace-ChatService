package grpc

import (
	"context"
	"time"

	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/room_invite/usecase"
	"github.com/MingPV/ChatService/pkg/apperror"
	roominvitepb "github.com/MingPV/ChatService/proto/room_invite"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcRoomInviteHandler struct {
	roomInviteUseCase usecase.RoomInviteUseCase
	roominvitepb.UnimplementedRoomInviteServiceServer
}

func NewGrpcRoomInviteHandler(uc usecase.RoomInviteUseCase) *GrpcRoomInviteHandler {
	return &GrpcRoomInviteHandler{roomInviteUseCase: uc}
}

// ------------------ Handlers ------------------

func (h *GrpcRoomInviteHandler) CreateRoomInvite(ctx context.Context, req *roominvitepb.CreateRoomInviteRequest) (*roominvitepb.CreateRoomInviteResponse, error) {
	invite := &entities.RoomInvite{
		RoomId:     uint(req.RoomId),
		Sender:     uuid.MustParse(req.Sender),
		InviteTo:   uuid.MustParse(req.InviteTo),
		IsAccepted: req.IsAccepted,
		IsDenied:   req.IsDenied,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.roomInviteUseCase.CreateRoomInvite(invite); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	return &roominvitepb.CreateRoomInviteResponse{
		Invite: toProtoRoomInvite(invite),
	}, nil
}

func (h *GrpcRoomInviteHandler) FindRoomInviteByID(ctx context.Context, req *roominvitepb.FindRoomInviteByIDRequest) (*roominvitepb.FindRoomInviteByIDResponse, error) {
	invite, err := h.roomInviteUseCase.FindByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &roominvitepb.FindRoomInviteByIDResponse{Invite: toProtoRoomInvite(invite)}, nil
}

func (h *GrpcRoomInviteHandler) FindAllRoomInvitesBySender(ctx context.Context, req *roominvitepb.FindAllRoomInvitesBySenderRequest) (*roominvitepb.FindAllRoomInvitesBySenderResponse, error) {
	invites, err := h.roomInviteUseCase.FindAllBySender(uuid.MustParse(req.Sender))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoInvites []*roominvitepb.RoomInvite
	for _, inv := range invites {
		protoInvites = append(protoInvites, toProtoRoomInvite(inv))
	}

	return &roominvitepb.FindAllRoomInvitesBySenderResponse{Invites: protoInvites}, nil
}

func (h *GrpcRoomInviteHandler) FindAllRoomInvitesByInviteTo(ctx context.Context, req *roominvitepb.FindAllRoomInvitesByInviteToRequest) (*roominvitepb.FindAllRoomInvitesByInviteToResponse, error) {
	invites, err := h.roomInviteUseCase.FindAllByInviteTo(uuid.MustParse(req.InviteTo))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoInvites []*roominvitepb.RoomInvite
	for _, inv := range invites {
		protoInvites = append(protoInvites, toProtoRoomInvite(inv))
	}

	return &roominvitepb.FindAllRoomInvitesByInviteToResponse{Invites: protoInvites}, nil
}

func (h *GrpcRoomInviteHandler) PatchRoomInvite(ctx context.Context, req *roominvitepb.PatchRoomInviteRequest) (*roominvitepb.PatchRoomInviteResponse, error) {
	updated := &entities.RoomInvite{
		ID:         uint(req.Id),
		RoomId:     uint(req.RoomId),
		Sender:     uuid.MustParse(req.Sender),
		InviteTo:   uuid.MustParse(req.InviteTo),
		IsAccepted: req.IsAccepted,
		IsDenied:   req.IsDenied,
		UpdatedAt:  time.Now(),
	}

	if err := h.roomInviteUseCase.PatchInvite(int(req.Id), updated); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	return &roominvitepb.PatchRoomInviteResponse{Invite: toProtoRoomInvite(updated)}, nil
}

func (h *GrpcRoomInviteHandler) DeleteRoomInvite(ctx context.Context, req *roominvitepb.DeleteRoomInviteRequest) (*roominvitepb.DeleteRoomInviteResponse, error) {
	if err := h.roomInviteUseCase.DeleteInvite(int(req.Id)); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &roominvitepb.DeleteRoomInviteResponse{Message: "room invite deleted"}, nil
}

// ------------------ Helpers ------------------

func toProtoRoomInvite(inv *entities.RoomInvite) *roominvitepb.RoomInvite {
	return &roominvitepb.RoomInvite{
		Id:         int32(inv.ID),
		RoomId:     int32(inv.RoomId),
		Sender:     inv.Sender.String(),
		InviteTo:   inv.InviteTo.String(),
		IsAccepted: inv.IsAccepted,
		IsDenied:   inv.IsDenied,
		CreatedAt:  timestamppb.New(inv.CreatedAt),
		UpdatedAt:  timestamppb.New(inv.UpdatedAt),
	}
}
