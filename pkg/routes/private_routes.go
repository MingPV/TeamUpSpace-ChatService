package routes

import (
	// middleware "github.com/MingPV/ChatService/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterPrivateRoutes(app fiber.Router, db *mongo.Database) {

	// route := app.Group("/api/v1", middleware.JWTMiddleware())

}
