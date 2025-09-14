package entities

import (
	"time"

	"github.com/google/uuid"
)

type Friend struct {
	ID    	  uint    `json:"id" bson:"_id,omitempty"`
    UserID    uuid.UUID `bson:"user_id" json:"user_id"`
    FriendID  uuid.UUID `bson:"friend_id" json:"friend_id"`
    Status    string    `bson:"status" json:"status"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
