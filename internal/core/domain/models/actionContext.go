package models

import (
	"kororo/internal/core/domain"
	"strings"
)

type ActionPipelineContext struct {
	actionsResponses []*ActionResponse
	fieldValues      []FieldValue
}

func (s *ActionPipelineContext) GetField(field *ActionField) (string, error) {
	for _, fieldValue := range s.fieldValues {
		if s.text(fieldValue.Field.Name) == s.text(field.Name) {
			return fieldValue.Value, nil
		}
	}

	for _, actionResponse := range s.actionsResponses {
		fieldValue := s.getFieldByActionResponse(actionResponse, field)
		if fieldValue != "" {
			return fieldValue, nil
		}
	}

	return "", domain.ErrFieldRequiredByAction
}

func (s *ActionPipelineContext) getFieldByActionResponse(actionResponse *ActionResponse, field *ActionField) string {

	for _, fieldResponse := range actionResponse.ResponseFields {
		if s.text(fieldResponse.Name) == s.text(field.Name) {
			return fieldResponse.Value
		}
	}

	return ""
}

func (s *ActionPipelineContext) text(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

func (s *ActionPipelineContext) GetAllFieldValuesActionResponses() []*ActionsResponseFields {
	var fields []*ActionsResponseFields

	for _, actionResponse := range s.actionsResponses {
		fields = append(fields, actionResponse.ResponseFields...)
	}

	return fields
}

func (s *ActionPipelineContext) GetAllActionResponses() []*ActionResponse {
	return s.actionsResponses
}

func NewActionPipelineContext(fieldValues []FieldValue) *ActionPipelineContext {
	return &ActionPipelineContext{
		fieldValues: fieldValues,
	}
}

func (s *ActionPipelineContext) AddActionResponse(actionResponse *ActionResponse) {
	s.actionsResponses = append(s.actionsResponses, actionResponse)
}
