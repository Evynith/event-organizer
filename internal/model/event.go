package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title             string             `json:"title"`
	Description_small string             `json:"description_small"`
	Description_large string             `json:"description_large"`
	DateOfEvent       primitive.DateTime `bson:"date" json:"date"`
	Organizer         primitive.ObjectID `bson:"organizer,omitempty" json:"organizer,omitempty"`
	Place             string             `json:"place"`
	Status            bool               `json:"status"` //0 eraser, 1 published
}

type Events []*Event
