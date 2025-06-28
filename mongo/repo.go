package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	DBName         string
	CollectionName string
}

func (r MongoRepo) getCollection() (*mongo.Collection, error) {
	client, err := Connect(context.Background())
	if err != nil {
		return nil, err
	}
	return client.Database(r.DBName).Collection(r.CollectionName), nil
}

func (r MongoRepo) InsertOne(doc interface{}) error {
	coll, err := r.getCollection()
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(context.Background(), doc)
	return err
}

func (r MongoRepo) FindOne(result interface{}, filter bson.M) error {
	coll, err := r.getCollection()
	if err != nil {
		return err
	}
	return coll.FindOne(context.Background(), filter).Decode(result)
}
