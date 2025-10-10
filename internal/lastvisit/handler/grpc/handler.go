package grpc

import (
	"context"
	"fmt"

	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/lastvisit/usecase"
	"github.com/MingPV/ChatService/pkg/apperror"
	lastvisitpb "github.com/MingPV/ChatService/proto/lastvisit"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcLastvisitHandler struct {
	lastvisitUseCase usecase.LastVisitUseCase
	lastvisitpb.UnimplementedLastvisitServiceServer
}

func NewGrpcLastvisitHandler(uc usecase.LastVisitUseCase) *GrpcLastvisitHandler {
	return &GrpcLastvisitHandler{lastvisitUseCase: uc}
}

// func (h *GrpcLastvisitHandler) CreateLastvisit(ctx context.Context, req *lastvisitpb.CreateLastvisitRequest) (*lastvisitpb.CreateLastvisitResponse, error){
// 	userUUID, err := uuid.Parse(req.UserId)
// 	if err != nil {
// 		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
// 	}

// 	lastvisit := &entities.Lastvisit{
//         UserID:    userUUID,
//         Lastvisit: time.Now(),
//     }

// 	err = h.lastvisitUseCase.CreateLastvisit(lastvisit)
//     if err != nil {
//         return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
//     }

// 	return &lastvisitpb.CreateLastvisitResponse{Lastvisit: toProtoLastvisit(lastvisit)}, nil
// }

func (h *GrpcLastvisitHandler) FindByUserID(ctx context.Context, req *lastvisitpb.FindByUserIDRequest) (*lastvisitpb.FindByUserIDResponse, error){
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	lastvisit, err := h.lastvisitUseCase.FindByUserID(userUUID, int(req.RoomId))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &lastvisitpb.FindByUserIDResponse{Lastvisit: toProtoLastvisit(lastvisit)}, nil
}

func (h *GrpcLastvisitHandler) UpdateLastvisit(ctx context.Context, req *lastvisitpb.UpdateLastvisitRequest) (*lastvisitpb.UpdateLastvisitResponse, error){
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	fmt.Println(req)
	
	updatedLastvisit, err := h.lastvisitUseCase.UpdateLastvisit(userUUID, int(req.RoomId))
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &lastvisitpb.UpdateLastvisitResponse{Lastvisit: toProtoLastvisit(updatedLastvisit)}, nil
}

func toProtoLastvisit(lvs *entities.Lastvisit) *lastvisitpb.Lastvisit {
	return &lastvisitpb.Lastvisit{
		UserId: lvs.UserID.String(),
		Lastvisit: timestamppb.New(lvs.Lastvisit),
		RoomId: int32(lvs.RoomID),
	}
}