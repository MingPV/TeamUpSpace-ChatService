package entities

import (
	"time"

	"github.com/google/uuid"
)
type Lastvisit struct {
	UserID 		uuid.UUID 	`bson:"user_id" json:"user_id"`
	Lastvisit 	time.Time 	`bson:"lastvisit" json:"lastvisit"`
	RoomID 		int 		`bson:"room_id" json:"room_id"`	
}