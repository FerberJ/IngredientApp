package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProvider struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var mongoProvider *MongoProvider

func NewMongoProvider(uri, dbName string) (*MongoProvider, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	setProvider(&MongoProvider{Client: client, Database: db})
	return mongoProvider, nil
}

func setProvider(provider *MongoProvider) {
	mongoProvider = provider
}

func GetProvider() *MongoProvider {
	return mongoProvider
}

func (mp *MongoProvider) Close() error {
	return mp.Client.Disconnect(context.TODO())
}
