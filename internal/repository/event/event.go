package repository

import (
	"context"
	"fmt"
	"time"

	"main/internal/database"
	model "main/internal/model"

	"go.mongodb.org/mongo-driver/bson"
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

func filterEvent(title string, date1 string, date2 string, state string, ids []primitive.ObjectID) bson.M {
	filter := bson.M{}
	nDate1, _ := time.Parse("2006-01-02", date1)
	nDate2, _ := time.Parse("2006-01-02", date2)
	nDate2 = nDate2.Add(24 * time.Hour)

	date1Time := primitive.NewDateTimeFromTime(nDate1)
	date2Time := primitive.NewDateTimeFromTime(nDate2)

	if title != "" {
		newSearch := bson.M{
			"$search": title,
		}
		filter["$text"] = newSearch
	}
	if len(ids) > 0 {
		newSearch := bson.M{
			"$in": ids,
		}
		filter["_id"] = newSearch
	}
	if state == "published" {
		filter["state"] = true
	} else if state == "eraser" {
		filter["state"] = false
	}
	if date1 != "" && date2 != "" {
		newSearch := bson.M{
			"$gte": date1Time,
			"$lte": date2Time,
		}
		filter["date"] = newSearch
	} else if date1 != "" && date2 == "" {
		newSearch := bson.M{
			"$gte": date1Time,
		}
		filter["date"] = newSearch
	} else if date1 == "" && date2 != "" {
		newSearch := bson.M{
			"$lte": date2Time,
		}
		filter["date"] = newSearch
	}
	fmt.Println(filter)
	return filter
}

func Read(title string, date1 string, date2 string, state string, user []primitive.ObjectID) (model.Events, error) {
	filter := filterEvent(title, date1, date2, state, user)
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
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
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
			"date":              event.DateOfEvent,
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
