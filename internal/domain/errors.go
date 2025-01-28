package domain

import "errors"

var (
	ErrInvalidFieldTypes        = errors.New("invalid field types")
	ErrFieldsNotUnique          = errors.New("fields not unique")
	ErrRequiredFieldsMissing    = errors.New("required fields missing")
	ErrInvalidFieldValues       = errors.New("invalid field values")
	ErrDuplicatedValue          = errors.New("duplicated value")
	ErrRequiredFieldsNotPresent = errors.New("required fields not present")
	ErrDuplicatedFields         = errors.New("duplicated fields")
)
