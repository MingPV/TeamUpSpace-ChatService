package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID    	 	uint    	`json:"id" bson:"_id,omitempty"`
	RoomId		uint 		`json:"room_id" bson:"room_id"`
	Message  	string		`json:"message" bson:"message"`
	Sender		uuid.UUID	`json:"sender" bson:"sender"`
	CreatedAt time.Time 	`json:"created_at" bson:"created_at"`
    UpdatedAt time.Time 	`json:"updated_at" bson:"updated_at"`
}