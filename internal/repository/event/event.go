package repository

import (
	"context"

	"main/internal/database"
	model "main/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = database.GetCollection("event")
var ctx = context.Background()

func Create(event model.Event) (interface{}, error) {
	var err error
	result, err := collection.InsertOne(ctx, event)
	if err != nil {
		return model.Event{}, err
	}
	return result.InsertedID, nil
}

func Read(filter primitive.M) (model.Events, error) {
	elems, err := collection.Find(ctx, filter)
	var events model.Events
	if err != nil {
		return nil, err
	}

	for elems.Next(ctx) {
		var event model.Event
		err = elems.Decode(&event)

		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func ReadOne(filter primitive.M) (model.Event, error) {
	var err error
	var result model.Event
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return model.Event{}, err
	}
	return result, nil

}

func Update(update primitive.M, filter primitive.M) error {
	var err error
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func Delete(filter primitive.M) error {
	var err error
	_, err = collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil
}
