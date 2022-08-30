package repository

import (
	"context"

	"main/internal/database"
	model "main/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = database.GetCollection("event")
var ctx = context.Background()

func Create(event model.Event) error {
	var err error
	_, err = collection.InsertOne(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

func Read() (model.Events, error) {
	filter := bson.M{}
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

func ReadOne(eventId string) (model.Event, error) {
	var err error
	oid, _ := primitive.ObjectIDFromHex(eventId)
	filter := bson.M{"_id": oid}

	var result model.Event
	event := collection.FindOne(ctx, filter).Decode(&result)

	if event != nil {
		return model.Event{}, err
	}
	return result, nil

}

func Update(event model.Event, eventId string) error {
	var err error
	oid, _ := primitive.ObjectIDFromHex(eventId)
	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"title":             event.Title,
			"description_small": event.Description_small,
			"description_large": event.Description_large,
			"date":              event.Date,
			"Organizer":         event.Organizer,
			"Place":             event.Place,
			"Status":            event.Status,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func Delete(eventId string) error {
	var err error
	oid, _ := primitive.ObjectIDFromHex(eventId)
	filter := bson.M{"_id": oid}

	_, err = collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil
}
