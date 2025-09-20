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

type MongoFriendRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoFriendRepository(db *mongo.Database) FriendRepository {
	return &MongoFriendRepository{
		db:   db,
		coll: db.Collection("friends"),
	}
}

type friendDoc struct {
	ID    	  int    	`bson:"_id,omitempty"`
    UserID    uuid.UUID `bson:"user_id"`
    FriendID  uuid.UUID `bson:"friend_id"`
    Status    string    `bson:"status"`
    CreatedAt time.Time `bson:"created_at"`
    UpdatedAt time.Time `bson:"updated_at"`
}

type counterDoc struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

func (r *MongoFriendRepository) getNextSequence(ctx context.Context, name string) (int, error) {
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

func (r *MongoFriendRepository)	Save(friend *entities.Friend) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nextID, err := r.getNextSequence(ctx, "friends")
	if err != nil {
		return err
	}

	_, err = r.coll.InsertOne(ctx, friendDoc{
		ID: nextID,
		UserID: friend.UserID,
		FriendID: friend.FriendID,
		Status: friend.Status,
		CreatedAt: friend.CreatedAt,
		UpdatedAt: friend.UpdatedAt,
	})

	if err != nil {
		return err
	}

	friend.ID = uint(nextID)
	return nil
}

func (r *MongoFriendRepository)	FindAll() ([]*entities.Friend, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.Friend
	for cur.Next(ctx) {
		var d friendDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, &entities.Friend{
			ID: uint(d.ID),
			UserID: d.UserID,
			FriendID: d.FriendID,
			Status: d.Status,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
func (r *MongoFriendRepository)	FindAllByUserId(userId uuid.UUID) ([]*entities.Friend, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"$or": []bson.M{
        {"user_id": userId},
        {"friend_id": userId},
    },}
	cur, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.Friend
	for cur.Next(ctx) {
		var d friendDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, &entities.Friend{
			ID: uint(d.ID),
			UserID: d.UserID,
			FriendID: d.FriendID,
			Status: d.Status,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoFriendRepository)	FindAllByIsFriend(userId uuid.UUID) ([]*entities.Friend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"status": "friend",
		"$or": []bson.M{
			{"user_id": userId},
			{"friend_id": userId},
		},
	}
	cur, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.Friend
	for cur.Next(ctx) {
		var d friendDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, &entities.Friend{
			ID: uint(d.ID),
			UserID: d.UserID,
			FriendID: d.FriendID,
			Status: d.Status,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoFriendRepository) IsMyfriend(userId uuid.UUID, friendId uuid.UUID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{
				"user_id":   userId,
				"friend_id": friendId,
			},
			{
				"user_id":   friendId,
				"friend_id": userId,
			},
		},
	}

	var friend *entities.Friend
	err := r.coll.FindOne(ctx, filter).Decode(&friend)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "not friend", nil
		}
		return "", err // real error
	}

	if friend.Status == "friend" {
		return "friend", nil
	}

	if friend.Status == "pending" {
		if friend.FriendID == userId {
			return "asked", nil
		}
		if friend.UserID == userId {
			return "pending", nil
		}
	}

	return "not friend", nil //fallback
}

func (r *MongoFriendRepository) FindAllFriendRequests(userId uuid.UUID) ([]*entities.Friend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"status": "pending", "friend_id" : userId}
	cur, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.Friend
	for cur.Next(ctx) {
		var d friendDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, &entities.Friend{
			ID: uint(d.ID),
			UserID: d.FriendID,
			FriendID: d.UserID,
			Status: d.Status,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoFriendRepository) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

func (r *MongoFriendRepository) Update(userId uuid.UUID, friendId uuid.UUID) (*entities.Friend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":   friendId,
		"friend_id": userId,
	}

	update := bson.M{
		"$set": bson.M{"status": "friend"},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedFriend entities.Friend
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedFriend)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &updatedFriend, nil
}

func (r *MongoFriendRepository) FindByID(id int) (*entities.Friend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	var f friendDoc
	err := r.coll.FindOne(ctx, filter).Decode(&f)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &entities.Friend{}, err
	}
	if err != nil {
		return nil, err
	}
	return &entities.Friend{
		ID: uint(f.ID),
		UserID: f.FriendID,
		FriendID: f.UserID,
		Status: f.Status,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}, nil
}


