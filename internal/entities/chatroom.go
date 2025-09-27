package entities

import (
	"time"

	"github.com/google/uuid"
)

type Chatroom struct {
	ID    	  	uint    	`json:"id" bson:"_id,omitempty"`
    RoomName    string 	    `bson:"room_name" json:"room_name"`
    IsGroup  	bool      	`bson:"is_group" json:"is_group"`
    Owner       uuid.UUID   `bson:"owner" json:"owner"`
    CreatedAt 	time.Time 	`bson:"created_at" json:"created_at"`
    UpdatedAt 	time.Time 	`bson:"updated_at" json:"updated_at"`
}