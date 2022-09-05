package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	usr      = "one"
	pwd      = "pass"
	host     = "mongo"
	port     = 27017
	database = "organization"
)

func GetCollection(collection string) *mongo.Collection {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", usr, pwd, host, port, database)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	inspectError(err)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	inspectError(err)

	return client.Database(database).Collection(collection)
}

func inspectError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
