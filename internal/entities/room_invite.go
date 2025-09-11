package entities

import (
	"time"

	"github.com/google/uuid"
)

type RoomInvite struct {
	ID    	  	uint    	`json:"id" bson:"_id,omitempty"`
	RoomId	  	uint 		`json:"room_id" bson:"room_id"`
	Sender		uuid.UUID	`json:"sender" bson:"sender"`
	InviteTo	uuid.UUID	`json:"invite_to" bson:"invite_to"`
	IsAccepted	bool		`json:"is_accepted" bson:"is_accepted"`
	IsDenied	bool		`json:"is_denied" bson:"is_denied"`
	CreatedAt 	time.Time 	`json:"created_at" bson:"created_at"`
    UpdatedAt 	time.Time 	`json:"updated_at" bson:"updated_at"`
}