package repository

import (
	"github.com/google/uuid"

	"context"
	"errors"
	"time"

	"github.com/MingPV/ChatService/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoLastvisitRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoLastvisitRepository(db *mongo.Database) LastvisitRepository {
	return &MongoLastvisitRepository{
		db:   db,
		coll: db.Collection("lastvisits"),
	}
}

type lastvisitDoc struct {
	UserID uuid.UUID `bson:"user_id"`
	Lastvisit time.Time `bson:"lastvisit"`
	RoomID int `bson:"room_id"`
}


// func (r *MongoLastvisitRepository) Save(lastvisit *entities.Lastvisit) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	_, err := r.coll.InsertOne(ctx, lastvisitDoc{
//         UserID:    lastvisit.UserID,
//         Lastvisit: lastvisit.Lastvisit,
// 		RoomID: lastvisit.RoomID,
//     })
//     if err != nil {
//         return err
//     }

// 	return nil
// }

func (r *MongoLastvisitRepository) FindByUserId(userId uuid.UUID, roomId int) (*entities.Lastvisit, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var doc lastvisitDoc
    err := r.coll.FindOne(ctx, bson.M{"user_id": userId, "room_id" : roomId}).Decode(&doc)
    if err != nil {
       if errors.Is(err, mongo.ErrNoDocuments) {
            return &entities.Lastvisit{}, nil // ไม่มี record
        }
        return nil, err
    }

    return &entities.Lastvisit{
        UserID:    doc.UserID,
        Lastvisit: doc.Lastvisit,
		RoomID: doc.RoomID,
    }, nil
}

func (r *MongoLastvisitRepository) Patch(userId uuid.UUID, roomId int) (*entities.Lastvisit, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    now := time.Now()

    // Update lastVisit หรือ insert ถ้าไม่มี
    opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
    var updatedDoc lastvisitDoc
    err := r.coll.FindOneAndUpdate(
        ctx,
        bson.M{"user_id": userId, "room_id": roomId},
        bson.M{"$set": bson.M{"lastvisit": now}},
        opts,
    ).Decode(&updatedDoc)
    if err != nil {
        return nil, err
    }

    return &entities.Lastvisit{
        UserID:    updatedDoc.UserID,
        Lastvisit: updatedDoc.Lastvisit,
		RoomID: updatedDoc.RoomID,
    }, nil
}
