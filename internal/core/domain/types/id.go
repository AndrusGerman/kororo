package types

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Id bson.ObjectID

func NewId() Id {
	return Id(bson.NewObjectID())
}

func NewIdByString(hex string) (Id, error) {
	var id, err = bson.ObjectIDFromHex(hex)
	return Id(id), err
}
