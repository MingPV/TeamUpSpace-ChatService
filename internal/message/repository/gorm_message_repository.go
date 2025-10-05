package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/MingPV/ChatService/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMessageRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoMessageRepository(db *mongo.Database) MessageRepository {
	return &MongoMessageRepository{
		db:   db,
		coll: db.Collection("messages"),
	}
}

type messageDoc struct {
	ID    	 	int    	`bson:"_id,omitempty"`
	RoomId		uint 		`bson:"room_id"`
	Message  	string		`bson:"message"`
	Sender		uuid.UUID	`bson:"sender"`
	CreatedAt 	time.Time 	`bson:"created_at"`
    UpdatedAt 	time.Time 	`bson:"updated_at"`
}

type counterDoc struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

func (r *MongoMessageRepository) getNextSequence(ctx context.Context, name string) (int, error) {
	counters := r.db.Collection("counters")
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var out counterDoc
	err := counters.FindOneAndUpdate(
		ctx,
		bson.M{"_id": name},
		bson.M{"$inc": bson.M{"seq": 1}},
		opts,
	).Decode(&out)

	if errors.Is(err, mongo.ErrNoDocuments) {
		_, ierr := counters.InsertOne(ctx, counterDoc{ID: name, Seq: 1})
		if ierr != nil {
			return 0, ierr
		}
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	if out.Seq == 0 {
		return 1, nil
	}
	return out.Seq, nil
}

func (r *MongoMessageRepository) Save(message *entities.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel();

	nextID, err := r.getNextSequence(ctx, "messages")
	if err != nil {
		return err
	}

	_, err = r.coll.InsertOne(ctx, messageDoc{
		ID: nextID,
		RoomId: message.RoomId,
		Message: message.Message,
		Sender: message.Sender,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	})

	if err != nil {
		return err
	}

	message.ID = uint(nextID)
	return nil
}

func (r *MongoMessageRepository) FindAllByRoomID(roomId int) ([]*entities.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"room_id": roomId}
	cur, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.Message
	for cur.Next(ctx) {
		var m messageDoc
		if err := cur.Decode(&m); err != nil {
			return nil, err
		}
		results = append(results, &entities.Message{
			ID: uint(m.ID),
			RoomId: m.RoomId,
			Message: m.Message,
			Sender: m.Sender,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoMessageRepository) DeleteAllMessagesByRoomID(roomId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"room_id" : roomId,
	}

	_, err := r.coll.DeleteMany(ctx, filter)
	return err
}

func (r *MongoMessageRepository) FindByRoomId(roomId int) (*entities.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"room_id": roomId}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})

	var message entities.Message

	err := r.coll.FindOne(ctx, filter, opts).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // no messages found
		}
		return nil, err
	}

	return &message, nil

}