package domain

import "errors"

var (
	ErrFieldsRequired = errors.New("fields are required")

	ErrIntentionNotFound = errors.New("intention not found")

	ErrMultipleIntentionSend = errors.New("multiple intention send error")

	ErrFieldDescriptionNotFound = errors.New("field description not found")

	ErrActionTypeNotSupported = errors.New("action type not supported")

	ErrFieldRequiredByAction = errors.New("field required by action")

	ErrClearJson = errors.New("error clear json")
)
