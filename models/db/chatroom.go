package models

import "github.com/kamva/mgm/v3"

type ChatRoom struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string   `json:"name" bson:"name"`
	Type             string   `json:"type" bson:"type"` // "group" or "private"
	Members          []string `json:"members" bson:"members"`
}

// CollectionName returns the name of the collection that stores ChatRoom documents.
func (model *ChatRoom) CollectionName() string {
	return "chatrooms"
}
