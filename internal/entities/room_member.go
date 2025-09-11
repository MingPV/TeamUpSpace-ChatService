package entities

import (
	"time"

	"github.com/google/uuid"
)

type RoomMember struct {
	ID    	  	uint    	`json:"id" bson:"_id,omitempty"`
	RoomId		uint 		`json:"room_id" bson:"room_id"`
	UserId 		uuid.UUID	`json:"user_id" bson:"user_id"`
	CreatedAt 	time.Time 	`json:"created_at" bson:"created_at"`
    UpdatedAt 	time.Time 	`json:"updated_at" bson:"updated_at"`
}