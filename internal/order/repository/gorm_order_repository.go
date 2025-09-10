package repository

import (
	"context"
	"errors"
	"time"

	"github.com/MingPV/ChatService/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NOTE: Replaced GORM repository with MongoDB implementation in the same package and file.

type MongoOrderRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoOrderRepository(db *mongo.Database) OrderRepository {
	return &MongoOrderRepository{
		db:   db,
		coll: db.Collection("orders"),
	}
}

type orderDoc struct {
	ID    int     `bson:"_id,omitempty"`
	Total float64 `bson:"total"`
}

type counterDoc struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

func (r *MongoOrderRepository) Save(order *entities.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nextID, err := r.getNextSequence(ctx, "orders")
	if err != nil {
		return err
	}

	_, err = r.coll.InsertOne(ctx, orderDoc{
		ID:    nextID,
		Total: order.Total,
	})
	if err != nil {
		return err
	}
	order.ID = uint(nextID)
	return nil
}

func (r *MongoOrderRepository) FindAll() ([]*entities.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*entities.Order
	for cur.Next(ctx) {
		var d orderDoc
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		results = append(results, &entities.Order{
			ID:    uint(d.ID),
			Total: d.Total,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoOrderRepository) FindByID(id int) (*entities.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var d orderDoc
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &entities.Order{}, err
	}
	if err != nil {
		return nil, err
	}
	return &entities.Order{
		ID:    uint(d.ID),
		Total: d.Total,
	}, nil
}

func (r *MongoOrderRepository) Patch(id int, order *entities.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{}
	update["total"] = order.Total

	_, err := r.coll.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *MongoOrderRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoOrderRepository) getNextSequence(ctx context.Context, name string) (int, error) {
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
