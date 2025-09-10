package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	// Order
	orderHandler "github.com/MingPV/ChatService/internal/order/handler/rest"
	orderRepository "github.com/MingPV/ChatService/internal/order/repository"
	orderUseCase "github.com/MingPV/ChatService/internal/order/usecase"
)

func RegisterPublicRoutes(app fiber.Router, db *mongo.Database) {

	api := app.Group("/api/v1")

	// === Dependency Wiring ===

	// Dependency wiring for Orders using MongoDB
	orderRepo := orderRepository.NewMongoOrderRepository(db)
	orderService := orderUseCase.NewOrderService(orderRepo)
	orderHandler := orderHandler.NewHttpOrderHandler(orderService)

	// === Public Routes ===

	// Order routes
	orderGroup := api.Group("/orders")
	orderGroup.Get("/", orderHandler.FindAllOrders)
	orderGroup.Get("/:id", orderHandler.FindOrderByID)
	orderGroup.Post("/", orderHandler.CreateOrder)
	orderGroup.Patch("/:id", orderHandler.PatchOrder)
	orderGroup.Delete("/:id", orderHandler.DeleteOrder)
}
