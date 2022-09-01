package model

type Filter struct {
	Title     string `json:"title"`
	DateSince string `bson:"since" json:"since"`
	DateUntil string `bson:"until" json:"until"`
	Status    string `json:"status"`
}
