package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSchemaHandler_Create(t *testing.T) {
	ctx := context.Background()
	_, humaApi := humatest.New(t)

	path := "/schema"

	InitSchemaRoutes(humaApi, NewSchemaHandler(NewSchemaServiceMock(NewSchemaRepositoryMock())))

	_ = t.Run("success", func(t *testing.T) {
		// arrange
		var reqBody = `
		{
			"fields": [
				{
					"name": "email",
					"type": "string",
					"required": true,
					"unique": true
				},
				{
					"name": "phone",
					"type": "integer",
					"required": true,
					"unique": true
				},
				{
					"name": "name",
					"type": "integer"
				}
			]
		}
		`
		// act
		res := humaApi.PostCtx(ctx, path, bytes.NewBufferString(reqBody))

		// assert
		if assert.Equal(t, http.StatusCreated, res.Code) {
			var resBody struct {
				ID string `json:"id"`
			}
			_ = json.Unmarshal(res.Body.Bytes(), &resBody)
			_ = assert.Equal(t, "67696ff2e3f76ec9d8e8dc3b", resBody.ID)
		}
	})

	_ = t.Run("invalid field type", func(t *testing.T) {
		// arrange
		var reqBody = `
		{
			"fields": [
				{
					"name": "email",
					"type": "string",
					"required": true,
					"unique": true
				},
				{
					"name": "phone",
					"type": "ineger",
					"required": true,
					"unique": true
				}
			]
		}
		`

		// act
		res := humaApi.PostCtx(ctx, path, bytes.NewBufferString(reqBody))

		// assert
		if assert.Equal(t, http.StatusBadRequest, res.Code) {
			var body huma.ErrorModel
			_ = json.Unmarshal(res.Body.Bytes(), &body)
			_ = assert.Equal(t, "Bad Request", body.Title)
			_ = assert.Equal(t, http.StatusBadRequest, body.Status)
			_ = assert.Equal(t, "invalid field types", body.Detail)
		}
	})

	_ = t.Run("fields not unique", func(t *testing.T) {
		// arrange
		var reqBody = `
		{
			"fields": [
				{
					"name": "email",
					"type": "string",
					"required": true,
					"unique": true
				},
				{
					"name": "email",
					"type": "string",
					"required": true,
					"unique": true
				}
			]
		}
		`

		// act
		res := humaApi.PostCtx(ctx, path, bytes.NewBufferString(reqBody))

		// assert
		if assert.Equal(t, http.StatusBadRequest, res.Code) {
			var body huma.ErrorModel
			_ = json.Unmarshal(res.Body.Bytes(), &body)
			_ = assert.Equal(t, "Bad Request", body.Title)
			_ = assert.Equal(t, http.StatusBadRequest, body.Status)
			_ = assert.Equal(t, "fields not unique", body.Detail)
		}
	})

	_ = t.Run("required fields not present", func(t *testing.T) {
		// arrange
		var reqBody = `
		{
			"fields": [
				{
					"name": "email",
					"type": "string",
					"required": true,
					"unique": true
				},
				{
					"name": "last_name",
					"type": "string",
					"required": true,
					"unique": true
				}
			]
		}
		`

		// act
		res := humaApi.PostCtx(ctx, path, bytes.NewBufferString(reqBody))

		// assert
		if assert.Equal(t, http.StatusBadRequest, res.Code) {
			var body huma.ErrorModel
			_ = json.Unmarshal(res.Body.Bytes(), &body)
			_ = assert.Equal(t, "Bad Request", body.Title)
			_ = assert.Equal(t, http.StatusBadRequest, body.Status)
			_ = assert.Equal(t, "required fields not present", body.Detail)
		}
	})

	_ = t.Run("without body", func(t *testing.T) {
		// act
		res := humaApi.PostCtx(ctx, path)

		// assert
		if assert.Equal(t, http.StatusBadRequest, res.Code) {
			var body huma.ErrorModel
			_ = json.Unmarshal(res.Body.Bytes(), &body)
			_ = assert.Equal(t, "Bad Request", body.Title)
			_ = assert.Equal(t, http.StatusBadRequest, body.Status)
			_ = assert.Equal(t, "request body is required", body.Detail)
		}
	})

	_ = t.Run("invalid body - fields", func(t *testing.T) {
		// arrange
		var reqBody = `
		{
			"error": [
				
			]
		}
		`

		// act
		res := humaApi.PostCtx(ctx, path, bytes.NewBufferString(reqBody))

		// assert
		if assert.Equal(t, http.StatusUnprocessableEntity, res.Code) {
			var body huma.ErrorModel
			_ = json.Unmarshal(res.Body.Bytes(), &body)
			_ = assert.Equal(t, "Unprocessable Entity", body.Title)
			_ = assert.Equal(t, http.StatusUnprocessableEntity, body.Status)
			_ = assert.Equal(t, "validation failed", body.Detail)
		}
	})

	_ = t.Run("invalid body - required name", func(t *testing.T) {
		// arrange
		var reqBody = `
		{
			"fields": [
				{
					"type": "string",
					"required": true,
				}
			]
		}
		`

		// act
		res := humaApi.PostCtx(ctx, path, bytes.NewBufferString(reqBody))

		// assert
		if assert.Equal(t, http.StatusUnprocessableEntity, res.Code) {
			var body huma.ErrorModel
			_ = json.Unmarshal(res.Body.Bytes(), &body)
			_ = assert.Equal(t, "Unprocessable Entity", body.Title)
			_ = assert.Equal(t, http.StatusUnprocessableEntity, body.Status)
			_ = assert.Equal(t, "validation failed", body.Detail)
		}
	})

	_ = t.Run("invalid body - required type", func(t *testing.T) {
		// arrange
		var reqBody = `
		{
			"fields": [
				{
					"name": "email",
				}
			]
		}
		`

		// act
		res := humaApi.PostCtx(ctx, path, bytes.NewBufferString(reqBody))

		// assert
		if assert.Equal(t, http.StatusUnprocessableEntity, res.Code) {
			var body huma.ErrorModel
			_ = json.Unmarshal(res.Body.Bytes(), &body)
			_ = assert.Equal(t, "Unprocessable Entity", body.Title)
			_ = assert.Equal(t, http.StatusUnprocessableEntity, body.Status)
			_ = assert.Equal(t, "validation failed", body.Detail)
		}
	})
}
