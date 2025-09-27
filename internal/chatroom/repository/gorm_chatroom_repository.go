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

type MongoChatroomRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoChatroomRepository(db *mongo.Database) ChatroomRepository {
	return &MongoChatroomRepository{
		db:   db,
		coll: db.Collection("chatrooms"),
	}
}

type chatroomDoc struct {
	ID 			int    		`bson:"_id,omitempty"`
	RoomName	string		`bson:"room_name"`
	IsGroup		bool 		`bson:"is_group"`
	Owner 		uuid.UUID  `bson:"owner" json:"owner"`
	CreatedAt 	time.Time 	`bson:"created_at"`
    UpdatedAt 	time.Time 	`bson:"updated_at"`
}

type counterDoc struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

func (r *MongoChatroomRepository) getNextSequence(ctx context.Context, name string) (int, error) {
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

func (r *MongoChatroomRepository)	Save(chatroom *entities.Chatroom) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nextID, err := r.getNextSequence(ctx, "chatrooms")
	if err != nil {
		return err
	}

	_, err = r.coll.InsertOne(ctx, chatroomDoc{
		ID: nextID,
		RoomName: chatroom.RoomName,
		IsGroup: chatroom.IsGroup,
		Owner: chatroom.Owner,
		CreatedAt: chatroom.CreatedAt,
		UpdatedAt: chatroom.UpdatedAt,
	})

	if err != nil {
		return err
	}

	chatroom.ID = uint(nextID)
	return nil
}

func (r *MongoChatroomRepository)	Patch(id int, chatroom *entities.Chatroom) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{}
	update["room_name"] = chatroom.RoomName
	update["owner"] = chatroom.Owner
	update["is_group"] = chatroom.IsGroup

	_, err := r.coll.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *MongoChatroomRepository)	FindByID(id int) (*entities.Chatroom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ch chatroomDoc
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&ch)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &entities.Chatroom{}, err
	}
	if err != nil {
		return nil, err
	}
	return &entities.Chatroom{
		ID:    uint(ch.ID),
		RoomName: ch.RoomName,
		IsGroup: ch.IsGroup,
		Owner: ch.Owner,
		CreatedAt: ch.CreatedAt,
		UpdatedAt: ch.UpdatedAt,
	}, nil
}

func (r *MongoChatroomRepository)	Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}