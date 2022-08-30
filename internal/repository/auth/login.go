package auth_repository

import (
	"context"
	"main/internal/database"
	"main/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

var collection = database.GetCollection("users")
var ctx = context.Background()

func SearchUser(username string) (model.User, error) {
	var err error
	filter := bson.M{"username": username}

	var result model.User
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return model.User{}, err
	}
	return result, nil

}

func PersistToken(username string, token string) error {
	var err error
	filter := bson.M{"username": username}
	newT := bson.M{
		"$set": bson.M{"token": token},
	}

	_, err = collection.UpdateOne(ctx, filter, newT)
	if err != nil {
		return err
	}
	return nil
}

func ExistsToken(token string) bool {
	var err error
	filter := bson.M{"token": token}

	var result model.User
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return false
	}
	return true
}
