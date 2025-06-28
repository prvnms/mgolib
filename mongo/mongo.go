package mongo

import (
	"context"
	"errors"
	"sync"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	clientInstance *mongo.Client
	once           sync.Once
)

func Connect(ctx context.Context) (*mongo.Client, error) {
	var err error
	once.Do(func() {
		uri := viper.GetString("MongoURL")
		if uri == "" {
			err = errors.New("MongoURL not found in config")
			return
		}
		clientInstance, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			return
		}
		err = clientInstance.Ping(ctx, readpref.Primary())
	})
	return clientInstance, err
}
