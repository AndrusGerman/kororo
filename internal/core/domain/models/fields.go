package models

import (
	"kororo/internal/core/domain/types"
)

type Field struct {
	*Base       `bson:",inline" json:",inline"`
	Description string          `bson:"description" json:"description"`
	Name        string          `bson:"name" json:"name"`
	Type        types.FieldType `bson:"type" json:"type"`
}

func NewField(description string, fieldType types.FieldType) *Field {
	return &Field{
		Base: NewBase(),
		//Description: description,
		Type: fieldType,
	}
}
