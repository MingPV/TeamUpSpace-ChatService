package routes

import (
	userHandler "github.com/MingPV/ChatService/internal/user/handler/rest"
	userRepository "github.com/MingPV/ChatService/internal/user/repository"
	userUseCase "github.com/MingPV/ChatService/internal/user/usecase"
	middleware "github.com/MingPV/ChatService/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterPrivateRoutes(app fiber.Router, db *gorm.DB) {

	route := app.Group("/api/v1", middleware.JWTMiddleware())

	userRepo := userRepository.NewGormUserRepository(db)
	ChatService := userUseCase.NewChatService(userRepo)
	userHandler := userHandler.NewHttpUserHandler(ChatService)

	route.Get("/me", userHandler.GetUser)

}
