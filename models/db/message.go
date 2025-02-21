package models

import "github.com/kamva/mgm/v3"

type Message struct {
	mgm.DefaultModel `bson:",inline"`
	From             string `json:"from" bson:"from"`
	Text             string `json:"text" bson:"text"`
	RoomID           string `json:"room_id" bson:"room_id"`
}

func (model *Message) CollectionName() string {
	return "messages"
}
