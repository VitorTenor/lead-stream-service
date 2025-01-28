package domain

import (
	"mime/multipart"
	"strconv"
)

type File struct {
	SchemaId string
	File     *multipart.FileHeader
}

func ValidateDuplicatedFields(headers []string) bool {
	seen := make(map[string]struct{})
	for _, field := range headers {
		if _, ok := seen[field]; ok {
			return false
		}
		seen[field] = struct{}{}
	}

	return true
}

func ValidateRequiredFields(headers []string) bool {
	requiredFields := map[string]string{
		"phone": "",
		"email": "",
	}

	seen := make(map[string]struct{})
	for _, field := range headers {
		seen[field] = struct{}{}
	}

	isValid := true
	for k, _ := range requiredFields {
		if _, ok := seen[k]; !ok {
			isValid = false
			break
		}
	}

	return isValid
}

func ValidateRequiredFieldsFromSchema(headers []string, sf []SchemaField) bool {
	seen := make(map[string]struct{})
	for _, field := range headers {
		seen[field] = struct{}{}
	}

	for _, header := range sf {
		if header.Required {
			if _, ok := seen[header.Name]; !ok {
				return false
			}
		}
	}

	return true
}

func ValueFromType(value string, t string) (interface{}, error) {
	switch {
	case t == "string":
		return value, nil
	case t == "float":
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	case t == "integer":
		i, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return i, nil
	case t == "boolean":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		return b, nil
	case t == "date":
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return i, nil
	case t == "time":
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return i, nil
	case t == "datetime":
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return i, nil
	default:
		return nil, ErrInvalidFieldValues
	}
}
