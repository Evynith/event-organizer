package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"main/internal/database"
	model "main/internal/model"
)

var collection = database.GetCollection("inscription")
var ctx = context.Background()

func Create(inscription model.Inscription) error {
	var err error
	_, err = collection.InsertOne(ctx, inscription)
	if err != nil {
		return err
	}
	return nil
}

func Read(idUser primitive.ObjectID) (model.Inscriptions, error) {
	filter := bson.M{"user": idUser}
	elems, err := collection.Find(ctx, filter)
	var inscriptions model.Inscriptions
	if err != nil {
		return nil, err
	}

	for elems.Next(ctx) {
		var inscription model.Inscription
		err = elems.Decode(&inscription)

		if err != nil {
			return nil, err
		}
		inscriptions = append(inscriptions, &inscription)
	}

	return inscriptions, nil
}
