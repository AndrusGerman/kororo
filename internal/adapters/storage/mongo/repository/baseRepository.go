package repository

import (
	"context"

	mongodb "kororo/internal/adapters/storage/mongo"
	"kororo/internal/core/domain/types"
	"kororo/internal/core/ports"

	"github.com/AndrusGerman/go-criteria"
	criteriatomongodb "github.com/AndrusGerman/go-criteria/driver/criteria-to-mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func newBaseRepository[T types.IBase](mongoService *mongodb.Mongo, collectionName string) ports.BaseRepository[T] {
	return &BaseRepository[T]{
		mongoService:   mongoService,
		collectionName: collectionName,
	}
}

type BaseRepository[T any] struct {
	mongoService   *mongodb.Mongo
	collectionName string
}

func (br *BaseRepository[T]) GetById(ctx context.Context, Id types.Id) (*T, error) {
	var base = new(T)
	var collection = br.Collection()
	var result = collection.FindOne(ctx, bson.M{"_id": bson.ObjectID(Id)})
	var err = result.Decode(base)
	return base, err
}

func (br *BaseRepository[T]) Delete(ctx context.Context, Id types.Id) error {
	var collection = br.Collection()
	var _, err = collection.DeleteOne(ctx, bson.M{"_id": bson.ObjectID(Id)})
	return err
}

func (br *BaseRepository[T]) Search(ctx context.Context, filter criteria.Criteria) ([]*T, error) {

	var vfilter = criteriatomongodb.NewCriteriaToMongodb().Convert(nil, filter, nil)

	var elements = new([]*T)
	var err error
	var collection = br.Collection()
	var cur *mongo.Cursor

	if cur, err = collection.Aggregate(ctx, vfilter); err != nil {
		return nil, err
	}
	if err = cur.All(ctx, elements); err != nil {
		return nil, err
	}

	return *elements, nil
}

func (ctx *BaseRepository[T]) Collection() *mongo.Collection {
	return ctx.mongoService.Collection(ctx.collectionName)
}

func (br *BaseRepository[T]) Update(ctx context.Context, set *T) error {
	var _, err = br.Collection().UpdateOne(ctx, bson.M{"_id": "id"}, bson.D{
		{
			Key: "$set",

			Value: set,
		},
	})
	return err
}

func (br *BaseRepository[T]) Create(ctx context.Context, element *T) error {
	var _, err = br.Collection().InsertOne(ctx, element)
	return err
}
