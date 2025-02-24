package repository

import (
	"context"
	"gotth/template/backend/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository struct {
	Provider   *db.MongoProvider
	Collection string
}

type Repository interface {
	InsertDocument(document any, ctx *context.Context) (primitive.ObjectID, error)
	FindDocument(filter bson.M, ctx *context.Context) (bson.M, error)
	FindDocuments(filter bson.M, ctx *context.Context) ([]bson.M, error)
	UpdateDocument(filter bson.M, update bson.M, ctx *context.Context) error
	FindDocumentsFields(filter bson.M, fields bson.M, ctx *context.Context) ([]bson.M, error)
}

func NewBaseRepository(provider *db.MongoProvider, collection string) *BaseRepository {
	return &BaseRepository{
		Provider:   provider,
		Collection: collection,
	}
}

func checkContext(ctx *context.Context) *context.Context {
	if ctx == nil {
		defaultCtx := context.Background()
		ctx = &defaultCtx
	}

	return ctx
}

func (r *BaseRepository) getCollection() *mongo.Collection {
	collection := r.Provider.Database.Collection(r.Collection)
	return collection
}

func (r *BaseRepository) InsertDocument(document any, ctx *context.Context) (primitive.ObjectID, error) {
	ctx = checkContext(ctx)
	collection := r.getCollection()
	result, err := collection.InsertOne(*ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *BaseRepository) FindDocument(filter bson.M, ctx *context.Context) (bson.M, error) {
	var result bson.M
	ctx = checkContext(ctx)
	collection := r.getCollection()
	err := collection.FindOne(*ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *BaseRepository) FindDocuments(filter bson.M, ctx *context.Context) ([]bson.M, error) {
	var results []bson.M
	ctx = checkContext(ctx)
	collection := r.getCollection()
	cursor, err := collection.Find(*ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(*ctx)

	for cursor.Next(*ctx) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			return nil, err
		}
		results = append(results, document)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *BaseRepository) FindDocumentsFields(filter bson.M, fields bson.M, ctx *context.Context) ([]bson.M, error) {
	var results []bson.M
	ctx = checkContext(ctx)
	collection := r.getCollection()

	findOptions := options.Find().SetProjection(fields)

	cursor, err := collection.Find(*ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(*ctx)

	for cursor.Next(*ctx) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			return nil, err
		}
		results = append(results, document)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *BaseRepository) UpdateDocument(filter bson.M, update bson.M, ctx *context.Context) error {
	ctx = checkContext(ctx)
	collection := r.getCollection()
	_, err := collection.UpdateOne(*ctx, filter, bson.M{"$set": update})
	return err
}
