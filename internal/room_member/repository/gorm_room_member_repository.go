package repository

import (
	"context"
	"errors"
	"time"

	"github.com/MingPV/ChatService/internal/entities"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRoomMemberRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoRoomMemberRepository(db *mongo.Database) RoomMemberRepository {
	return &MongoRoomMemberRepository{
		db:   db,
		coll: db.Collection("room_members"),
	}
}

type roomMemberDoc struct {
	ID        int       `bson:"_id,omitempty"`
	RoomId    uint      `bson:"room_id"`
	UserId    uuid.UUID `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type counterDoc struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

// Save multiple users into a room
func (r *MongoRoomMemberRepository) Save(roomId uint, userIDs []uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	var docs []interface{}

	for _, userID := range userIDs {
		nextID, err := r.getNextSequence(ctx, "room_members")
		if err != nil {
			return err
		}
		docs = append(docs, roomMemberDoc{
			ID:        nextID,
			RoomId:    roomId,
			UserId:    userID,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	_, err := r.coll.InsertMany(ctx, docs)
	return err
}

// FindAllByRoomID returns all members of a room
func (r *MongoRoomMemberRepository) FindAllByRoomID(roomId uint) ([]*entities.RoomMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := r.coll.Find(ctx, bson.M{"room_id": roomId})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.RoomMember
	for cur.Next(ctx) {
		var d roomMemberDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, &entities.RoomMember{
			ID:        uint(d.ID),
			RoomId:    d.RoomId,
			UserId:    d.UserId,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoRoomMemberRepository) FindAllByUserID(userId uuid.UUID) ([]*entities.RoomMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userId}}},
		{{
			Key: "$lookup", Value: bson.M{
				"from":         "chatrooms", // collection ของห้อง
				"localField":   "room_id",   // room_id ใน room_members
				"foreignField": "_id",       // _id ใน chatrooms
				"as":           "room",      // output field ที่ map เข้า struct
			},
		}},
		{{
			Key: "$unwind", Value: bson.M{
				"path":                       "$room",
				"preserveNullAndEmptyArrays": true,
			},
		}},
	}

	// ต้องใช้ Aggregate แทน Find
	cur, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.RoomMember
	for cur.Next(ctx) {
		var d entities.RoomMember
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, &d)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}


// FindAllByRoomIDAndUserID returns a single member in a room
func (r *MongoRoomMemberRepository) FindAllByRoomIDAndUserID(roomId uint, userId uuid.UUID) (*entities.RoomMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var d roomMemberDoc
	err := r.coll.FindOne(ctx, bson.M{"room_id": roomId, "user_id": userId}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entities.RoomMember{
		ID:        uint(d.ID),
		RoomId:    d.RoomId,
		UserId:    d.UserId,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}, nil
}

// DeleteByRoomIDAndUserID deletes a specific room member
func (r *MongoRoomMemberRepository) DeleteByRoomIDAndUserID(roomId uint, userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.DeleteOne(ctx, bson.M{
		"room_id": roomId,
		"user_id": userId,
	})
	return err
}

// DeleteAllByRoomID deletes all members of a room
func (r *MongoRoomMemberRepository) DeleteAllByRoomID(roomId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.DeleteMany(ctx, bson.M{"room_id": roomId})
	return err
}

// getNextSequence generates auto-increment ID
func (r *MongoRoomMemberRepository) getNextSequence(ctx context.Context, name string) (int, error) {
	counters := r.db.Collection("counters")
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

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

func (r *MongoRoomMemberRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
