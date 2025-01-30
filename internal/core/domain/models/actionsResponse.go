package models

import "kororo/internal/core/domain/types"

type ActionResponse struct {
	ActionId       types.Id                 `bson:"action_id" json:"action_id"`
	Response       string                   `bson:"response" json:"response"`
	ResponseFields []*ActionsResponseFields `bson:"response_fields" json:"response_fields"`
	Status         string
}

type ActionsResponseFields struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
