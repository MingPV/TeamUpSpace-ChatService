package app

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"

	GrpcOrderHandler "github.com/MingPV/ChatService/internal/order/handler/grpc"
	orderRepository "github.com/MingPV/ChatService/internal/order/repository"
	orderUseCase "github.com/MingPV/ChatService/internal/order/usecase"
	orderpb "github.com/MingPV/ChatService/proto/order"

	GrpcFriendHandler "github.com/MingPV/ChatService/internal/friend/handler/grpc"
	friendRepository "github.com/MingPV/ChatService/internal/friend/repository"
	friendUseCase "github.com/MingPV/ChatService/internal/friend/usecase"
	friendpb "github.com/MingPV/ChatService/proto/friend"

	GrpcMessageHandler "github.com/MingPV/ChatService/internal/message/handler/grpc"
	messageRepository "github.com/MingPV/ChatService/internal/message/repository"
	messageUseCase "github.com/MingPV/ChatService/internal/message/usecase"
	messagepb "github.com/MingPV/ChatService/proto/message"

	GrpcChatroomHandler "github.com/MingPV/ChatService/internal/chatroom/handler/grpc"
	chatroomRepository "github.com/MingPV/ChatService/internal/chatroom/repository"
	chatroomUseCase "github.com/MingPV/ChatService/internal/chatroom/usecase"
	chatroompb "github.com/MingPV/ChatService/proto/chatroom"

	GrpcRoomMemberHandler "github.com/MingPV/ChatService/internal/room_member/handler/grpc"
	roommemberRepository "github.com/MingPV/ChatService/internal/room_member/repository"
	roommemberUseCase "github.com/MingPV/ChatService/internal/room_member/usecase"
	roommemberpb "github.com/MingPV/ChatService/proto/room_member"

	GrpcRoomInviteHandler "github.com/MingPV/ChatService/internal/room_invite/handler/grpc"
	roominviteRepository "github.com/MingPV/ChatService/internal/room_invite/repository"
	roominviteUseCase "github.com/MingPV/ChatService/internal/room_invite/usecase"
	roominvitepb "github.com/MingPV/ChatService/proto/room_invite"

	"github.com/MingPV/ChatService/pkg/config"
	"github.com/MingPV/ChatService/pkg/database"
	"github.com/MingPV/ChatService/pkg/middleware"
	"github.com/MingPV/ChatService/pkg/routes"
)

// rest
func SetupRestServer(db *mongo.Database, cfg *config.Config) (*fiber.App, error) {
	app := fiber.New()
	middleware.FiberMiddleware(app)
	// comment out Swagger when testing
	// routes.SwaggerRoute(app)
	routes.RegisterPublicRoutes(app, db, cfg)
	routes.RegisterPrivateRoutes(app, db)
	routes.RegisterNotFoundRoute(app)
	return app, nil
}

// grpc
func SetupGrpcServer(db *mongo.Database, cfg *config.Config) (*grpc.Server, error) {
	s := grpc.NewServer()

	// Dependency wiring for Orders using MongoDB
	orderRepo := orderRepository.NewMongoOrderRepository(db)
	orderService := orderUseCase.NewOrderService(orderRepo)
	orderHandler := GrpcOrderHandler.NewGrpcOrderHandler(orderService)
	orderpb.RegisterOrderServiceServer(s, orderHandler)

	
	
	roommemberRepo := roommemberRepository.NewMongoRoomMemberRepository(db)
	roommemberService := roommemberUseCase.NewRoomMemberService(roommemberRepo)
	roommemberHandler := GrpcRoomMemberHandler.NewGrpcRoomMemberHandler(roommemberService)
	roommemberpb.RegisterRoomMemberServiceServer(s, roommemberHandler)
	
	roominviteRepo := roominviteRepository.NewMongoRoomInviteRepository(db)
	roominviteService := roominviteUseCase.NewRoomInviteService(roominviteRepo, roommemberRepo)
	roominviteHandler := GrpcRoomInviteHandler.NewGrpcRoomInviteHandler(roominviteService)
	roominvitepb.RegisterRoomInviteServiceServer(s, roominviteHandler)
	
	
	// Message streaming service
	msgRepo := messageRepository.NewMongoMessageRepository(db)
	msgUseCase := messageUseCase.NewMessageService(msgRepo)
	msgHandler := GrpcMessageHandler.NewGrpcMessageHandler(msgUseCase)
	messagepb.RegisterMessageServiceServer(s, msgHandler)
	
	chatroomRepo := chatroomRepository.NewMongoChatroomRepository(db)
	chatroomService := chatroomUseCase.NewChatroomService(chatroomRepo, roommemberRepo, msgRepo)
	chatroomHandler := GrpcChatroomHandler.NewGrpcChatroomHandler(chatroomService)
	chatroompb.RegisterChatroomServiceServer(s, chatroomHandler)
	
	friendRepo := friendRepository.NewMongoFriendRepository(db)
	friendService := friendUseCase.NewFriendService(friendRepo, chatroomRepo, roommemberRepo)
	friendHandler := GrpcFriendHandler.NewGrpcFriendHandler(friendService)
	friendpb.RegisterFriendServiceServer(s, friendHandler)

	return s, nil
}

// dependencies
func SetupDependencies(env string) (*mongo.Database, *config.Config, error) {
	cfg := config.LoadConfig(env)

	db, err := database.Connect(cfg.MongoURI, cfg.DBName)
	if err != nil {
		return nil, nil, err
	}

	return db, cfg, nil
}
