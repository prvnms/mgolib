package mongo

import (
	"bytes"
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	DBName string
	Client *mongo.Client
}

func (r *MongoRepo) getCollection(ctype interface{}) (*mongo.Collection, error) {
	cname := getType(ctype)
	return r.Client.Database(r.DBName).Collection(cname), nil
}

func (r *MongoRepo) InsertOne(doc interface{}) error {
	coll, err := r.getCollection(doc)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(context.Background(), doc)
	return err
}

func (r *MongoRepo) FindOne(result interface{}, filter bson.M) error {
	coll, err := r.getCollection(result)
	if err != nil {
		return err
	}
	return coll.FindOne(context.Background(), filter).Decode(result)
}

func getType(i interface{}) string {
	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	if isSlice(t) {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}
	return getFirstLowerCase(t.Name())
}

func isSlice(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Slice
}

func getFirstLowerCase(str string) string {
	bts := []byte(str)
	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]
	return string(bytes.Join([][]byte{lc, rest}, nil))
}

func isPointer(fld interface{}) bool {
	return reflect.TypeOf(fld).Kind() == reflect.Ptr
}
