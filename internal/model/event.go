package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	Title             string             `json:"title"`
	Description_small string             `json:"description_small"`
	Description_large string             `json:"description_large"`
	Date              primitive.DateTime `bson:"date" json:"date"`
	Organizer         primitive.ObjectID `bson:"organizer,omitempy" json:"organizer,omitempy"`
	Place             string             `json:"place"`
	Status            bool               `json:"status"` //0 eraser, 1 published
}

type Events []*Event
