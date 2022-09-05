package auth_repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"main/internal/database"
	"main/internal/model"
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

/*
Recibe un token codificado, usuario asociado y revisa en la base de datos si ha sido guardado anteriormente para dicho usuario
*/
func ExistsToken(token string, idUser string) bool {
	var err error
	oid, _ := primitive.ObjectIDFromHex(idUser)
	filter := bson.M{"token": token, "_id": oid}

	var result model.User
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil || result.Username == "" {
		return false
	}
	return true
}
