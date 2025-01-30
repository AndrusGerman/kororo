package domain

import "errors"

var (
	ErrFieldsRequired = errors.New("fields are required")

	ErrFieldDescriptionNotFound = errors.New("field description not found")

	ErrActionTypeNotSupported = errors.New("action type not supported")

	ErrFieldRequiredByAction = errors.New("field required by action")
)
