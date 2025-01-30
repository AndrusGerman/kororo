package models

type FieldValue struct {
	Field *Field `json:"field"`
	Value string `json:"value"`
}

func NewFieldValue(field *Field, value string) *FieldValue {
	return &FieldValue{
		Field: field,
		Value: value,
	}
}
