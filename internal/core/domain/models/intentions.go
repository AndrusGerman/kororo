package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Intention struct {
	*Base          `bson:",inline" json:",inline"`
	Description    string          `bson:"description" json:"description"`
	Fields         []*Field        `bson:"fields" json:"fields"`
	Actions        []bson.ObjectID `bson:"actions" json:"actions"`
	ResponseAction bson.ObjectID   `bson:"response_action" json:"response_action"`
}

func NewIntention(description string, fields []*Field) *Intention {
	return &Intention{
		Base:        NewBase(),
		Description: description,
		Fields:      fields,
	}
}
