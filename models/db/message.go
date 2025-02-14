package models

import "github.com/kamva/mgm/v3"

type Message struct {
	mgm.DefaultModel `bson:",inline"`
	RoomID           string `bson:"room_id" json:"room_id"`
	From             string `bson:"from" json:"from"`
	To               string `bson:"to" json:"to"`
	Content          string `bson:"content" json:"content"`
	Type             string `bson:"type" json:"type"`
}

func (model *Message) CollectionName() string {
	return "messages"
}
