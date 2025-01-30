package models

import (
	"kororo/internal/core/domain/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Base struct {
	Id        bson.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

func (b *Base) GetId() types.Id {
	return types.Id(b.Id)
}

func NewBase() *Base {
	var id = bson.NewObjectID()
	var now = time.Now()

	return &Base{
		Id:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
