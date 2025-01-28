package handlers

import (
	"errors"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"net/http"

	"github.com/vitortenor/lead-stream-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func handleError(err error) error {
	switch {
	case errors.Is(err, primitive.ErrInvalidHex),
		errors.Is(err, domain.ErrFieldsNotUnique),
		errors.Is(err, domain.ErrInvalidFieldTypes),
		errors.Is(err, domain.ErrInvalidFieldValues),
		errors.Is(err, domain.ErrDuplicatedValue),
		errors.Is(err, domain.ErrRequiredFieldsNotPresent),
		errors.Is(err, domain.ErrDuplicatedFields):
		return huma.NewError(http.StatusBadRequest, err.Error())

	case errors.Is(err, domain.ErrRequiredFieldsMissing):
		return huma.NewError(http.StatusBadRequest, err.Error())

	case errors.Is(err, mongo.ErrNoDocuments):
		return huma.NewError(http.StatusNotFound, err.Error())

	case mongo.IsDuplicateKeyError(err):
		return huma.NewError(http.StatusConflict, err.Error())

	default:
		return huma.NewError(http.StatusInternalServerError,
			fmt.Sprintf("Internal server error: %s", err.Error()))
	}
}
