package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/mongo"

	// Order
	orderHandler "github.com/MingPV/ChatService/internal/order/handler/rest"
	orderRepository "github.com/MingPV/ChatService/internal/order/repository"
	orderUseCase "github.com/MingPV/ChatService/internal/order/usecase"

	// Message gateway over gRPC
	messageGateway "github.com/MingPV/ChatService/internal/message/handler"
	messagepb "github.com/MingPV/ChatService/proto/message"
	"google.golang.org/grpc"

	"github.com/MingPV/ChatService/pkg/config"
)

func RegisterPublicRoutes(app fiber.Router, db *mongo.Database, cfg *config.Config) {

	api := app.Group("/api/v1")

	// === Dependency Wiring ===

	// Dependency wiring for Orders using MongoDB
	orderRepo := orderRepository.NewMongoOrderRepository(db)
	orderService := orderUseCase.NewOrderService(orderRepo)
	orderHandler := orderHandler.NewHttpOrderHandler(orderService)


	// WebSocket -> gRPC gateway client
	grpcConn, _ := grpc.Dial("localhost:"+cfg.GrpcPort, grpc.WithInsecure())
	msgClient := messagepb.NewMessageServiceClient(grpcConn)
	wsGateway := messageGateway.NewWebSocketGatewayHandler(msgClient)

	// === Public Routes ===

	// Order routes
	orderGroup := api.Group("/orders")
	orderGroup.Get("/", orderHandler.FindAllOrders)
	orderGroup.Get("/:id", orderHandler.FindOrderByID)
	orderGroup.Post("/", orderHandler.CreateOrder)
	orderGroup.Patch("/:id", orderHandler.PatchOrder)
	orderGroup.Delete("/:id", orderHandler.DeleteOrder)

	// Message websocket routes
	wsGroup := api.Group("/ws")
	wsGroup.Use("/rooms/:roomId", wsGateway.UpgradeMiddleware)
	wsGroup.Get("/rooms/:roomId", websocket.New(wsGateway.SubscribeRoomWebSocket))
}
