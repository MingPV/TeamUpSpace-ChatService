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

type MongoRoomInviteRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoRoomInviteRepository(db *mongo.Database) RoomInviteRepository {
	return &MongoRoomInviteRepository{
		db:   db,
		coll: db.Collection("room_invites"),
	}
}

// --- Mongo document mapping ---
type roomInviteDoc struct {
	ID        int       `bson:"_id,omitempty"`
	RoomId    uint      `bson:"room_id"`
	Sender    uuid.UUID `bson:"sender"`
	InviteTo  uuid.UUID `bson:"invite_to"`
	IsAccepted bool     `bson:"is_accepted"`
	IsDenied   bool     `bson:"is_denied"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type counterDoc struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

// --- CRUD operations ---
func (r *MongoRoomInviteRepository) Save(invite *entities.RoomInvite) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nextID, err := r.getNextSequence(ctx, "room_invites")
	if err != nil {
		return err
	}

	_, err = r.coll.InsertOne(ctx, roomInviteDoc{
		ID:         nextID,
		RoomId:     invite.RoomId,
		Sender:     invite.Sender,
		InviteTo:   invite.InviteTo,
		IsAccepted: invite.IsAccepted,
		IsDenied:   invite.IsDenied,
		CreatedAt:  invite.CreatedAt,
		UpdatedAt:  invite.UpdatedAt,
	})
	if err != nil {
		return err
	}

	invite.ID = uint(nextID)
	return nil
}

func (r *MongoRoomInviteRepository) FindByID(id int) (*entities.RoomInvite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var d roomInviteDoc
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &entities.RoomInvite{}, err
	}
	if err != nil {
		return nil, err
	}

	return r.toEntity(d), nil
}

func (r *MongoRoomInviteRepository) FindAllBySender(sender uuid.UUID) ([]*entities.RoomInvite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := r.coll.Find(ctx, bson.M{"sender": sender})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []*entities.RoomInvite{}, err
	}
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.RoomInvite
	for cur.Next(ctx) {
		var d roomInviteDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, r.toEntity(d))
	}
	return results, cur.Err()
}

func (r *MongoRoomInviteRepository) FindAllByRoomId(roomId int) ([]*entities.RoomInvite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := r.coll.Find(ctx, bson.M{"room_id": roomId})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []*entities.RoomInvite{}, err
	}
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.RoomInvite
	for cur.Next(ctx) {
		var d roomInviteDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, r.toEntity(d))
	}
	return results, cur.Err()
}


func (r *MongoRoomInviteRepository) FindAllByInviteTo(inviteTo uuid.UUID) ([]*entities.RoomInvite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := r.coll.Find(ctx, bson.M{"invite_to": inviteTo})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []*entities.RoomInvite{}, err
	}
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.RoomInvite
	for cur.Next(ctx) {
		var d roomInviteDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, r.toEntity(d))
	}
	return results, cur.Err()
}

func (r *MongoRoomInviteRepository) Patch(id int, invite *entities.RoomInvite) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{}
	if invite.RoomId != 0 {
		update["room_id"] = invite.RoomId
	}
	update["sender"] = invite.Sender
	update["invite_to"] = invite.InviteTo
	update["is_accepted"] = invite.IsAccepted
	update["is_denied"] = invite.IsDenied
	update["updated_at"] = invite.UpdatedAt

	_, err := r.coll.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *MongoRoomInviteRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// --- Helpers ---
func (r *MongoRoomInviteRepository) getNextSequence(ctx context.Context, name string) (int, error) {
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

func (r *MongoRoomInviteRepository) toEntity(d roomInviteDoc) *entities.RoomInvite {
	return &entities.RoomInvite{
		ID:         uint(d.ID),
		RoomId:     d.RoomId,
		Sender:     d.Sender,
		InviteTo:   d.InviteTo,
		IsAccepted: d.IsAccepted,
		IsDenied:   d.IsDenied,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}
