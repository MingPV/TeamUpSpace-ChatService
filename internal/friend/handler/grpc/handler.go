package grpc

import (
	"context"

	"github.com/MingPV/ChatService/internal/entities"
	"github.com/MingPV/ChatService/internal/friend/usecase"
	"github.com/MingPV/ChatService/pkg/apperror"
	friendpb "github.com/MingPV/ChatService/proto/friend"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcFriendHandler struct {
	friendUseCase usecase.FriendUseCase
	friendpb.UnimplementedFriendServiceServer
}

func NewGrpcFriendHandler(uc usecase.FriendUseCase) *GrpcFriendHandler {
	return &GrpcFriendHandler{friendUseCase: uc}
}

func (h *GrpcFriendHandler) CreateFriend(ctx context.Context, req *friendpb.CreateFriendRequest) (*friendpb.CreateFriendResponse, error) {
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	friendUUID, err := uuid.Parse(req.FriendId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	friend := &entities.Friend{
		UserID: userUUID,
		FriendID: friendUUID,
		Status: req.Status,
	}

	if err := h.friendUseCase.CreateFriend(friend); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &friendpb.CreateFriendResponse{Friend: toProtoFriend(friend)}, nil
}

func (h *GrpcFriendHandler) FindAllFriends(ctx context.Context, req *friendpb.FindAllFriendsRequest) (*friendpb.FindAllFriendsResponse, error) {
	friends, err := h.friendUseCase.FindAllFriends()
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoFriends []*friendpb.Friend
	for _, f := range friends {
		protoFriends = append(protoFriends, toProtoFriend(f))
	}

	return &friendpb.FindAllFriendsResponse{Friends: protoFriends}, nil
}

func (h *GrpcFriendHandler) FindAllFriendsByUserID(ctx context.Context, req *friendpb.FindAllFriendsByUserIDRequest) (*friendpb.FindAllFriendsByUserIDResponse, error) {
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	friends, err := h.friendUseCase.FindAllFriendsByUserID(userUUID)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoFriends []*friendpb.Friend
	for _, f := range friends {
		protoFriends = append(protoFriends, toProtoFriend(f))
	}

	return &friendpb.FindAllFriendsByUserIDResponse{Friends: protoFriends}, nil
}

func (h *GrpcFriendHandler) FindAllFriendsByIsFriend(ctx context.Context, req *friendpb.FindAllFriendsByIsFriendRequest) (*friendpb.FindAllFriendsByIsFriendResponse, error) {
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	friends, err := h.friendUseCase.FindAllFriendsByIsFriend(userUUID)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoFriends []*friendpb.Friend
	for _, f := range friends {
		protoFriends = append(protoFriends, toProtoFriend(f))
	}

	return &friendpb.FindAllFriendsByIsFriendResponse{Friends: protoFriends}, nil
}

func (h *GrpcFriendHandler) FindAllFriendRequests(ctx context.Context, req *friendpb.FindAllFriendRequestsRequest) (*friendpb.FindAllFriendRequestsResponse, error) {
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	friends, err := h.friendUseCase.FindAllFriendsRequests(userUUID)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	var protoFriends []*friendpb.Friend
	for _, f := range friends {
		protoFriends = append(protoFriends, toProtoFriend(f))
	}

	return &friendpb.FindAllFriendRequestsResponse{Friends: protoFriends}, nil
}

func (h *GrpcFriendHandler) DeleteFriend(ctx context.Context, req *friendpb.DeleteFriendRequest) (*friendpb.DeleteFriendResponse, error) {
	if err := h.friendUseCase.DeleteFriend(uint(req.Id)); err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	return &friendpb.DeleteFriendResponse{Message: "friend deleted"}, nil
}

func (h *GrpcFriendHandler) IsMyFriend(ctx context.Context, req *friendpb.IsMyFriendRequest) (*friendpb.IsMyFriendResponse, error) {
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	friendUUID, err := uuid.Parse(req.FriendId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	
	status, _ := h.friendUseCase.IsMyfriend(userUUID, friendUUID)
	return &friendpb.IsMyFriendResponse{Status: status}, nil
}

func (h *GrpcFriendHandler) AcceptFriend(ctx context.Context, req *friendpb.AcceptFriendRequest) (*friendpb.AcceptFriendResponse, error){
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}
	friendUUID, err := uuid.Parse(req.FriendId)
	if err != nil {
		return nil, status.Errorf(apperror.GRPCCode(err), "%s", err.Error())
	}

	friend, err := h.friendUseCase.AcceptFriend(userUUID, friendUUID)
	if err != nil {
		return nil, err
	}
	return &friendpb.AcceptFriendResponse{Friend: toProtoFriend(friend)}, nil
}


func toProtoFriend(f *entities.Friend) *friendpb.Friend {
	return &friendpb.Friend{
		Id:    int32(f.ID),
		UserId: f.UserID.String(),
		FriendId: f.FriendID.String(),
		Status: f.Status,
		CreatedAt: timestamppb.New(f.CreatedAt),
		UpdatedAt: timestamppb.New(f.CreatedAt),
	}
}


