package mongo

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MgoClientInstance *mongo.Client
	initOnce          sync.Once
	initErr           error
	dbName            string
)

type MongoConfig struct {
	URI                string
	DBName             string
	MaxPoolSize        uint64
	MinPoolSize        uint64
	MaxConnIdleTimeSec int
	ServerSelectionTO  int
	SocketTO           int
	ConnectTO          int
}

func Connect(ctx context.Context, config MongoConfig) (*mongo.Client, error) {
	initOnce.Do(func() {
		if config.URI == "" {
			initErr = errors.New("Mongo URI is empty")
			return
		}
		if config.DBName == "" {
			initErr = errors.New("Mongo DB name is empty")
			return
		}

		dbName = config.DBName

		clientOptions := options.Client().
			ApplyURI(config.URI).
			SetMaxPoolSize(config.MaxPoolSize).
			SetMinPoolSize(config.MinPoolSize).
			SetMaxConnIdleTime(time.Duration(config.MaxConnIdleTimeSec) * time.Second).
			SetServerSelectionTimeout(time.Duration(config.ServerSelectionTO) * time.Second).
			SetSocketTimeout(time.Duration(config.SocketTO) * time.Second).
			SetConnectTimeout(time.Duration(config.ConnectTO) * time.Second)

		MgoClientInstance, initErr = mongo.Connect(ctx, clientOptions)
		if initErr != nil {
			log.Println("Mongo connect error:", initErr)
			return
		}

		initErr = MgoClientInstance.Ping(ctx, readpref.Primary())
		if initErr != nil {
			log.Println("Mongo ping error:", initErr)
			return
		}

		log.Println("Mongo connected successfully")
	})

	if initErr != nil {
		MgoClientInstance = nil
		initOnce = sync.Once{}
	}

	return MgoClientInstance, initErr
}

func Disconnect(ctx context.Context) error {
	if MgoClientInstance != nil {
		err := MgoClientInstance.Disconnect(ctx)
		if err != nil {
			log.Println("Mongo disconnect error:", err)
			return err
		}
		log.Println("Mongo disconnected successfully")
	}
	return nil
}

func DBName() string {
	return dbName
}
