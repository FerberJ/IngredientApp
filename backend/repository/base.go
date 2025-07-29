package repository

import (
	"context"
	"encoding/json"
	"errors"
	"gotth/template/backend/db"

	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

type BaseRepository struct {
	Provider   *db.CloverProvider
	Collection string
}

type Repository interface {
	InsertDocument(document any, ctx *context.Context) (string, error)
	FindDocument(query *query.Query, ctx *context.Context) (map[string]any, error)
	FindDocuments(query *query.Query, ctx *context.Context) ([]map[string]any, error)
	UpdateDocument(query *query.Query, update map[string]any, ctx *context.Context) error
	//FindDocumentsFields(filter bson.M, fields bson.M, ctx *context.Context) ([]bson.M, error)
	DeleteDocument(query *query.Query, ctx *context.Context) error
}

func NewBaseRepository(provider *db.CloverProvider, collection string) *BaseRepository {
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

func (r *BaseRepository) getCollection() *c.DB {
	hasCollection, err := r.Provider.Database.HasCollection(r.Collection)
	if err != nil {
		return nil
	}

	if !hasCollection {
		err := r.Provider.Database.CreateCollection(r.Collection)
		if err != nil {
			return nil
		}
	}

	collection := r.Provider.Database
	return collection
}

func (r *BaseRepository) getDoc(document any) (*d.Document, error) {
	mapDoc := make(map[string]any)

	documentByte, err := json.Marshal(document)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(documentByte, &mapDoc)
	if err != nil {
		return nil, err
	}

	doc := d.NewDocument()
	doc.SetAll(mapDoc)

	return doc, nil
}

func (r *BaseRepository) InsertDocument(document any, ctx *context.Context) (string, error) {
	ctx = checkContext(ctx)
	collection := r.getCollection()
	doc, err := r.getDoc(document)
	if err != nil {
		return "", err
	}
	result, err := collection.InsertOne(r.Collection, doc)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *BaseRepository) FindDocument(query *query.Query, ctx *context.Context) (map[string]any, error) {
	ctx = checkContext(ctx)
	collection := r.getCollection()
	doc, err := collection.FindFirst(query)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, errors.New("document not found")
		// or return empty map, depending on your needs
	}
	return doc.AsMap(), nil
}

func (r *BaseRepository) FindDocuments(query *query.Query, ctx *context.Context) ([]map[string]any, error) {
	var results []map[string]any
	ctx = checkContext(ctx)
	collection := r.getCollection()
	docs, err := collection.FindAll(query)
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		results = append(results, doc.ToMap())
	}

	return results, nil
}

/*
func (r *BaseRepository) FindDocumentsFields(filter *query.Query, fields bson.M, ctx *context.Context) ([]bson.M, error) {
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
*/

func (r *BaseRepository) UpdateDocument(query *query.Query, update map[string]any, ctx *context.Context) error {
	ctx = checkContext(ctx)
	collection := r.getCollection()
	err := collection.Update(query, update)
	return err
}

func (r *BaseRepository) DeleteDocument(query *query.Query, ctx *context.Context) error {
	ctx = checkContext(ctx)
	collection := r.getCollection()
	err := collection.Delete(query)
	return err
}
