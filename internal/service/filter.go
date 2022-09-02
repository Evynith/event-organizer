package service

import (
	"main/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slices"
)

func AccessDraft(typeUser string) bool {
	access := []string{"admin"}
	rta := slices.Contains(access, typeUser)
	return rta
}

func CreateFilterID(id string) bson.M {
	oid, _ := primitive.ObjectIDFromHex(id)
	return bson.M{"_id": oid}
}

func CreateFilterListOfEvent(filter bson.M, inscriptionsList []primitive.ObjectID) bson.M {
	if len(inscriptionsList) > 0 {
		newSearch := bson.M{
			"$in": inscriptionsList,
		}
		filter["_id"] = newSearch
	}
	return filter
}

func CreateEventUpdate(event model.Event) bson.M {
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
	return update
}

func CreateFilterEvent(id string, draft bool) bson.M {
	filter := CreateFilterID(id)
	if !draft {
		filter["status"] = true
	}
	return filter
}

func CreateFilterEvents(filter model.Filter, ids []primitive.ObjectID, draft bool) bson.M {
	fltr := bson.M{}
	nDate1, _ := time.Parse("2006-01-02", filter.DateSince)
	nDate2, _ := time.Parse("2006-01-02", filter.DateUntil)
	nDate2 = nDate2.Add(24 * time.Hour)

	date1Time := primitive.NewDateTimeFromTime(nDate1)
	date2Time := primitive.NewDateTimeFromTime(nDate2)

	if filter.Title != "" {
		newSearch := bson.M{
			"$search": filter.Title,
		}
		fltr["$text"] = newSearch
	}
	if draft {
		if filter.Status == "published" {
			fltr["status"] = true
		} else if filter.Status == "draft" {
			fltr["status"] = false
		}
	} else {
		fltr["status"] = true
	}

	if filter.DateSince != "" && filter.DateUntil != "" {
		newSearch := bson.M{
			"$gte": date1Time,
			"$lte": date2Time,
		}
		fltr["date"] = newSearch
	} else if filter.DateSince != "" && filter.DateUntil == "" {
		newSearch := bson.M{
			"$gte": date1Time,
		}
		fltr["date"] = newSearch
	} else if filter.DateSince == "" && filter.DateUntil != "" {
		newSearch := bson.M{
			"$lte": date2Time,
		}
		fltr["date"] = newSearch
	}
	fltr = CreateFilterListOfEvent(fltr, ids)
	return fltr
}
