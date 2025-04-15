package db

import (
	"context"
	"gotth/template/backend/configuration"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProvider struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var mongoProvider *MongoProvider

func NewMongoProvider(cfg configuration.Configutration) (*MongoProvider, error) {
	clientOptions := options.Client().ApplyURI(cfg.MongoEndpoint)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.MongoDb)
	setMongoProvider(&MongoProvider{Client: client, Database: db})
	return mongoProvider, nil
}

func setMongoProvider(provider *MongoProvider) {
	mongoProvider = provider
}

func GetMongoProvider() *MongoProvider {
	return mongoProvider
}

func (mp *MongoProvider) Close() error {
	return mp.Client.Disconnect(context.TODO())
}
