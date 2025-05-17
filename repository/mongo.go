package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[T any] struct {
	Collection *mongo.Collection
}

func NewMongoRepository[T any](client *mongo.Client, dbName, collectionName string) *Repository[T] {
	coll := client.Database(dbName).Collection(collectionName)
	return &Repository[T]{Collection: coll}
}

func (r *Repository[T]) Save(item T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.InsertOne(ctx, item)
	return err
}

func (r *Repository[T]) FindByField(field string, value any) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{field: value}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var elem T
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	return results, nil
}
