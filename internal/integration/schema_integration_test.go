package integration

import (
	"bytes"
	"encoding/json"
	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestSchemaHandler_Create(t *testing.T) {
	srv, err := InitServerTest()
	if err != nil {
		t.Fatal(err)
	}

	schemaUrl := srv.URL + "/schema"

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
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, bytes.NewBufferString(reqBody))

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusCreated, res.StatusCode) {
				var resBody struct {
					ID string `json:"id"`
				}
				body, _ := io.ReadAll(res.Body)
				_ = json.Unmarshal(body, &resBody)
				_ = assert.NotEmpty(t, resBody.ID)
			}
		}
	})

	_ = t.Run("invalid body - invalid field type", func(t *testing.T) {
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
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, bytes.NewBufferString(reqBody))

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "invalid field types", body.Detail)
			}
		}
	})

	_ = t.Run("invalid body - fields not unique", func(t *testing.T) {
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
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, bytes.NewBufferString(reqBody))

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "fields not unique", body.Detail)
			}
		}
	})

	_ = t.Run("invalid body - required fields not present", func(t *testing.T) {
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
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, bytes.NewBufferString(reqBody))

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "required fields not present", body.Detail)
			}
		}
	})

	_ = t.Run("invalid body - without body", func(t *testing.T) {
		// act
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, nil)

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "request body is required", body.Detail)
			}
		}
	})

	_ = t.Run("invalid body - no fields", func(t *testing.T) {
		// arrange
		var reqBody = `
		{
			"error": [
	
			]
		}
		`

		// act
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, bytes.NewBufferString(reqBody))

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Unprocessable Entity", body.Title)
				_ = assert.Equal(t, http.StatusUnprocessableEntity, body.Status)
				_ = assert.Equal(t, "validation failed", body.Detail)
			}
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
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, bytes.NewBufferString(reqBody))

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "validation failed", body.Detail)
			}
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
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, bytes.NewBufferString(reqBody))

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "validation failed", body.Detail)
			}
		}
	})

	_ = t.Run("invalid body - created_at field", func(t *testing.T) {
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
				},
				{
					"name": "created_at",
					"type": "datetime"
				},
			]
		}
		`
		// act
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, bytes.NewBufferString(reqBody))

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "validation failed", body.Detail)
			}
		}
	})

	_ = t.Run("invalid body - updated_at field", func(t *testing.T) {
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
				},
				{
					"name": "updated_at",
					"type": "datetime"
				},
			]
		}
		`
		// act
		res, err := http.Post(schemaUrl, echo.MIMEApplicationJSON, bytes.NewBufferString(reqBody))

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "validation failed", body.Detail)
			}
		}
	})
}
