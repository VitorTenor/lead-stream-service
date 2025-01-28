package integration

import (
	"bytes"
	"encoding/json"
	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vitortenor/lead-stream-service/internal/tools"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileHandler_Upload(t *testing.T) {
	srv, err := InitServerTest()
	if err != nil {
		t.Fatal("Failed to initialize server:", err)
	}

	rootPath, err := tools.FindProjectRoot()
	if err != nil {
		t.Fatal("Failed to find project root:", err)
	}

	fileUrl := srv.URL + "/schema/{schemaId}/file"
	schemaId, invalidSchemaID := "67808a19c567c857d77d7f12", "67699e3d7887d75bb8523702"

	_ = t.Run("success", func(t *testing.T) {
		// arrange
		file, err := openFile(rootPath, "test_file_handler_success.csv")
		if err != nil {
			t.Fatal("Failed to open test file:", err)
		}
		defer file.Close()

		urlWithParams := strings.Replace(fileUrl, "{schemaId}", schemaId, 1)

		body, contentType, err := createMultipartForm(file)
		if err != nil {
			t.Fatal("Failed to create multipart form:", err)
		}

		// act
		res, err := makeRequest(urlWithParams, contentType, &body)
		if err != nil {
			t.Fatal("Failed to perform request:", err)
		}
		defer res.Body.Close()

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusOK, res.StatusCode) {
				var resBody struct {
					Message string `json:"message"`
				}
				body, _ := io.ReadAll(res.Body)
				_ = json.Unmarshal(body, &resBody)
				_ = assert.Equal(t, "File uploaded successfully", resBody.Message)
			}
		}
	})

	_ = t.Run("non existent schema", func(t *testing.T) {
		// arrange
		file, err := openFile(rootPath, "test_file_handler_success.csv")
		if err != nil {
			t.Fatal("Failed to open test file:", err)
		}
		defer file.Close()

		body, contentType, err := createMultipartForm(file)
		if err != nil {
			t.Fatal("Failed to create multipart form:", err)
		}

		urlWithParams := strings.Replace(fileUrl, "{schemaId}", invalidSchemaID, 1)

		// act
		res, err := makeRequest(urlWithParams, contentType, &body)
		if err != nil {
			t.Fatal("Failed to perform request:", err)
		}
		defer res.Body.Close()

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusNotFound, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Not Found", body.Title)
				_ = assert.Equal(t, http.StatusNotFound, body.Status)
				_ = assert.Equal(t, "mongo: no documents in result", body.Detail)
			}
		}
	})

	_ = t.Run("schema id is not a object id", func(t *testing.T) {
		// arrange
		file, err := openFile(rootPath, "test_file_handler_success.csv")
		if err != nil {
			t.Fatal("Failed to open test file:", err)
		}
		defer file.Close()

		body, contentType, err := createMultipartForm(file)
		if err != nil {
			t.Fatal("Failed to create multipart form:", err)
		}

		urlWithParams := strings.Replace(fileUrl, "{schemaId}", "123", 1)

		// act
		res, err := makeRequest(urlWithParams, contentType, &body)
		if err != nil {
			t.Fatal("Failed to perform request:", err)
		}
		defer res.Body.Close()

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "the provided hex string is not a valid ObjectID", body.Detail)
			}
		}
	})

	_ = t.Run("without body", func(t *testing.T) {
		// arrange
		urlWithParams := strings.Replace(fileUrl, "{schemaId}", "123", 1)

		// act
		res, err := http.Post(urlWithParams, echo.MIMEMultipartForm, nil)

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

	_ = t.Run("missing schema required field", func(t *testing.T) {
		// arrange
		file, err := openFile(rootPath, "test_file_handler_fail.csv")
		if err != nil {
			t.Fatal("Failed to open test file:", err)
		}
		defer file.Close()

		urlWithParams := strings.Replace(fileUrl, "{schemaId}", schemaId, 1)

		body, contentType, err := createMultipartForm(file)
		if err != nil {
			t.Fatal("Failed to create multipart form:", err)
		}

		// act
		res, err := makeRequest(urlWithParams, contentType, &body)
		if err != nil {
			t.Fatal("Failed to perform request:", err)
		}
		defer res.Body.Close()

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "required fields missing", body.Detail)
			}
		}
	})

	_ = t.Run("missing system required field", func(t *testing.T) {
		// arrange
		file, err := openFile(rootPath, "test_file_handler_fail_2.csv")
		if err != nil {
			t.Fatal("Failed to open test file:", err)
		}
		defer file.Close()

		urlWithParams := strings.Replace(fileUrl, "{schemaId}", schemaId, 1)

		body, contentType, err := createMultipartForm(file)
		if err != nil {
			t.Fatal("Failed to create multipart form:", err)
		}

		// act
		res, err := makeRequest(urlWithParams, contentType, &body)
		if err != nil {
			t.Fatal("Failed to perform request:", err)
		}
		defer res.Body.Close()

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "required fields missing", body.Detail)
			}
		}
	})

	_ = t.Run("repeated values", func(t *testing.T) {
		// arrange
		file, err := openFile(rootPath, "test_file_handler_fail_3.csv")
		if err != nil {
			t.Fatal("Failed to open test file:", err)
		}
		defer file.Close()

		urlWithParams := strings.Replace(fileUrl, "{schemaId}", schemaId, 1)

		body, contentType, err := createMultipartForm(file)
		if err != nil {
			t.Fatal("Failed to create multipart form:", err)
		}

		// act
		res, err := makeRequest(urlWithParams, contentType, &body)
		if err != nil {
			t.Fatal("Failed to perform request:", err)
		}
		defer res.Body.Close()

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "duplicated value", body.Detail)
			}
		}
	})

	_ = t.Run("repeated header values", func(t *testing.T) {
		// arrange
		file, err := openFile(rootPath, "test_file_handler_fail_4.csv")
		if err != nil {
			t.Fatal("Failed to open test file:", err)
		}
		defer file.Close()

		urlWithParams := strings.Replace(fileUrl, "{schemaId}", schemaId, 1)

		body, contentType, err := createMultipartForm(file)
		if err != nil {
			t.Fatal("Failed to create multipart form:", err)
		}

		// act
		res, err := makeRequest(urlWithParams, contentType, &body)
		if err != nil {
			t.Fatal("Failed to perform request:", err)
		}
		defer res.Body.Close()

		// assert
		if assert.NoError(t, err) {
			if assert.Equal(t, http.StatusBadRequest, res.StatusCode) {
				var body huma.ErrorModel
				_ = json.NewDecoder(res.Body).Decode(&body)
				_ = assert.Equal(t, "Bad Request", body.Title)
				_ = assert.Equal(t, http.StatusBadRequest, body.Status)
				_ = assert.Equal(t, "duplicated fields", body.Detail)
			}
		}
	})
}

func openFile(rootPath, fileName string) (*os.File, error) {
	return os.Open(filepath.Join(rootPath, "internal", "integration", "resources", "file", fileName))
}

func createMultipartForm(file *os.File) (bytes.Buffer, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return body, "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return body, "", err
	}
	writer.Close()
	return body, writer.FormDataContentType(), nil
}

func makeRequest(url, contentType string, body *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(echo.HeaderContentType, contentType)

	client := &http.Client{}
	return client.Do(req)
}
