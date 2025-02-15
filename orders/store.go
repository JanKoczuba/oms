package main

import (
	"context"
	pb "github.com/JanKoczuba/commons/api"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
)

const (
	DbName   = "orders"
	CollName = "orders"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client,
) *store {
	return &store{db}
}

func (s *store) Create(ctx context.Context, o Order) (bson.ObjectID, error) {
	col := s.db.Database(DbName).Collection(CollName)

	newOrder, err := col.InsertOne(ctx, o)

	id := newOrder.InsertedID.(bson.ObjectID)
	return id, err
}

func (s *store) Get(ctx context.Context, id, customerID string) (*Order, error) {
	col := s.db.Database(DbName).Collection(CollName)

	oID, _ := bson.ObjectIDFromHex(id)

	var o Order
	err := col.FindOne(ctx, bson.M{
		"_id":        oID,
		"customerID": customerID,
	}).Decode(&o)

	return &o, err
}

func (s *store) Update(ctx context.Context, id string, newOrder *pb.Order) error {
	col := s.db.Database(DbName).Collection(CollName)

	oID, _ := bson.ObjectIDFromHex(id)

	log.Printf("oid: %s", oID)
	log.Printf("id: %s", id)

	_, err := col.UpdateOne(ctx,
		bson.M{"_id": oID},
		bson.M{"$set": bson.M{
			"paymentLink": newOrder.PaymentLink,
			"status":      newOrder.Status,
		}})

	return err
}
