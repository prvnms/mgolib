package mongo

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func InitDefaultPID(dbName, collectionName string) {
	repo := &MongoRepo{
		DBName:         dbName,
		CollectionName: collectionName,
	}

	var seq PID
	err := repo.FindOne(&seq, bson.M{"key": "EXE_ID"})
	if err != nil {
		entry := PID{Key: "EXE_ID", Seq: 0}
		err = repo.InsertOne(entry)
		if err != nil {
			log.Println("DB init error:", err)
		}
	}
}
