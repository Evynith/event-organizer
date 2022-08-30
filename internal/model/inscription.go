package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Inscription struct {
	Event primitive.ObjectID `bson:"event" json:"event"`
	User  primitive.ObjectID `bson:"user" json:"user"`
}

type Inscriptions []*Inscription
